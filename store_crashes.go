package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Crash struct {
	Id             int
	CreatedOn      time.Time
	App            *App
	ProgramVersion *string
	IpAddrInternal *string
	CrashingLine   *string
	Sha1           [20]byte
}

type App struct {
	Name                   string
	Crashes                []*Crash
	PerDayCrashes          map[string][]*Crash
	PerCrashingLineCrashes map[string][]*Crash
}

// CrashesCount returns number of crashes for an app
func (a *App) CrashesCount() int {
	return len(a.Crashes)
}

// StoreCrashes is a store for crashes
type StoreCrashes struct {
	sync.Mutex
	dataDir       string
	crashes       []*Crash
	apps          []*App
	versions      []*string
	ips           map[string]*string
	crashingLines map[string]*string
	dataFile      *os.File
}

func (c *Crash) IpAddress() string {
	return ipAddrInternalToOriginal(*c.IpAddrInternal)
}

// CreatedOnDay returns a string version of created on
func (c *Crash) CreatedOnDay() string {
	return c.CreatedOn.Format("2006-01-02")
}

// CrashingLineCount returns number of unique crashing lines
func (c *Crash) CrashingLineCount() int {
	return len(c.App.PerCrashingLineCrashes[*c.CrashingLine])
}

// CrashesByCreatedOn is for sorting crashing by CreatedOn time
type CrashesByCreatedOn []*Crash

func (s CrashesByCreatedOn) Len() int {
	return len(s)
}
func (s CrashesByCreatedOn) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s CrashesByCreatedOn) Less(i, j int) bool {
	t1 := s[i].CreatedOn
	t2 := s[j].CreatedOn
	return t1.After(t2)
}

// GetAppByName returns App by its name
func (s *StoreCrashes) GetAppByName(appName string) *App {
	for _, app := range s.apps {
		if appName == app.Name {
			return app
		}
	}
	return nil
}

// FindOrCreateApp returns existing App or creates a new one for a given name
func (s *StoreCrashes) FindOrCreateApp(appName string) *App {
	if app := s.GetAppByName(appName); app != nil {
		return app
	}

	app := &App{
		Name:                   appName,
		Crashes:                make([]*Crash, 0),
		PerDayCrashes:          make(map[string][]*Crash),
		PerCrashingLineCrashes: make(map[string][]*Crash),
	}
	s.apps = append(s.apps, app)
	return app
}

func (s *StoreCrashes) FindOrCreateVersion(ver string) *string {
	for _, v := range s.versions {
		if *v == ver {
			return v
		}
	}
	s.versions = append(s.versions, &ver)
	return &ver
}

func (s *StoreCrashes) FindCrashingLine(str string) *string {
	if s2, ok := s.crashingLines[str]; ok {
		return s2
	}
	return nil
}

func (s *StoreCrashes) FindOrCreateCrashingLine(str string) *string {
	if s2, ok := s.crashingLines[str]; ok {
		return s2
	}
	s.crashingLines[str] = &str
	return &str
}

func (s *StoreCrashes) FindIp(str string) *string {
	if s2, ok := s.ips[str]; ok {
		return s2
	}
	return nil
}

func (s *StoreCrashes) FindOrCreateIp(str string) *string {
	if s2, ok := s.ips[str]; ok {
		return s2
	}
	s.ips[str] = &str
	return &str
}

func ipAddrInternalToOriginal(s string) string {
	// check if ipv4 in hex form
	if len(s) == 8 {
		if d, err := hex.DecodeString(s); err != nil {
			return s
		} else {
			return fmt.Sprintf("%d.%d.%d.%d", d[0], d[1], d[2], d[3])
		}
	}
	// other format (ipv6?)
	return s
}

func ipAddrToInternal(ipAddr string) string {
	var nums [4]uint32
	parts := strings.Split(ipAddr, ".")
	if len(parts) != 4 {
		// assuming it's ip v6
		return ipAddr
	}

	for n, p := range parts {
		num, _ := strconv.Atoi(p)
		nums[n] = uint32(num)
	}
	n := (nums[0] << 24) | (nums[1] << 16) + (nums[2] << 8) | nums[3]
	return fmt.Sprintf("%x", n)
}

// parse:
// C/vs1mJI02u0HBsHPceGfxy/Q+JE|1351741403|SumatraPDF|2.1.1|6e8e602f
func (s *StoreCrashes) parseCrash(line []byte) {
	parts := strings.Split(string(line[1:]), "|")
	if len(parts) != 6 {
		panic("len(parts) != 6")
	}
	msgSha1b64 := parts[0] + "="
	createdOnSecondsStr := parts[1]
	programName := parts[2]
	programVersion := parts[3]
	ipAddrInternal := parts[4]
	crashingLineTmp := parts[5]

	createdOnSeconds, err := strconv.Atoi(createdOnSecondsStr)
	if err != nil {
		panic("createdOnSeconds not a number")
	}
	createdOn := time.Unix(int64(createdOnSeconds), 0)

	msgSha1, err := base64.StdEncoding.DecodeString(msgSha1b64)
	if err != nil {
		panic("msgSha1b64 not valid base64")
	}
	if len(msgSha1) != 20 {
		panic("len(msgSha1) != 20")
	}

	programVersionInterned := s.FindOrCreateVersion(programVersion)
	ipAddr := s.FindOrCreateIp(ipAddrInternal)
	app := s.FindOrCreateApp(programName)
	crashingLine := s.FindOrCreateCrashingLine(crashingLineTmp)
	c := &Crash{
		Id:             len(s.crashes),
		App:            app,
		CreatedOn:      createdOn,
		ProgramVersion: programVersionInterned,
		IpAddrInternal: ipAddr,
		CrashingLine:   crashingLine,
	}
	copy(c.Sha1[:], msgSha1)

	if !s.MessageFileExists(c.Sha1[:]) {
		panic("message file doesn't exist")
	}
	s.appendCrash(c)
}

func (s *StoreCrashes) appendCrash(c *Crash) {
	s.crashes = append(s.crashes, c)
	c.App.Crashes = append(c.App.Crashes, c)
	day := c.CreatedOnDay()
	perDay, ok := c.App.PerDayCrashes[day]
	if !ok {
		perDay = make([]*Crash, 0)
	}
	perDay = append(perDay, c)
	c.App.PerDayCrashes[day] = perDay

	cl := *c.CrashingLine
	perCrashingLine, ok := c.App.PerCrashingLineCrashes[cl]
	if !ok {
		perCrashingLine = make([]*Crash, 0)
	}
	perCrashingLine = append(perCrashingLine, c)
	c.App.PerCrashingLineCrashes[cl] = perCrashingLine
}

func (s *StoreCrashes) readExistingCrashesData(fileDataPath string) error {
	d, err := ioutil.ReadFile(fileDataPath)
	if err != nil {
		return err
	}

	for len(d) > 0 {
		idx := bytes.IndexByte(d, '\n')
		if -1 == idx {
			// TODO: this could happen if the last record was only
			// partially written. Should I just ignore it?
			panic("idx shouldn't be -1")
		}
		line := d[:idx]
		d = d[idx+1:]
		c := line[0]
		if c == 'C' {
			s.parseCrash(line)
		} else {
			fmt.Printf("%q\n", string(line))
			panic("Unexpected line type")
		}
	}
	return nil
}

// NewStoreCrashes returns new StoreCrashes
func NewStoreCrashes(dataDir string) (*StoreCrashes, error) {
	dataFilePath := filepath.Join(dataDir, "data", "crashesdata.txt")
	CreateDirForFileMust(dataFilePath)
	store := &StoreCrashes{
		dataDir:       dataDir,
		crashes:       make([]*Crash, 0),
		apps:          make([]*App, 0),
		versions:      make([]*string, 0),
		ips:           make(map[string]*string),
		crashingLines: make(map[string]*string),
	}

	var err error
	if PathExists(dataFilePath) {
		err = store.readExistingCrashesData(dataFilePath)
		if err != nil {
			logger.Errorf("NewStoreCrashes(): readExistingCrashesData() failed with %s\n", err)
			return nil, err
		}
	} else {
		f, err := os.Create(dataFilePath)
		if err != nil {
			logger.Errorf("NewStoreCrashes(): os.Create(%s) failed with %s", dataFilePath, err)
			return nil, err
		}
		f.Close()
	}
	store.dataFile, err = os.OpenFile(dataFilePath, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		logger.Errorf("NewStoreCrashes(): os.OpenFile(%s) failed with %s", dataFilePath, err)
		return nil, err
	}
	logger.Noticef("crashes: %d, versions: %d, ips: %d, crashing lines: %d", len(store.crashes), len(store.versions), len(store.ips), len(store.crashingLines))

	return store, nil
}

// CrashesCount returns total number of crashes
func (s *StoreCrashes) CrashesCount() int {
	s.Lock()
	defer s.Unlock()
	return len(s.crashes)
}

func blobCrahesPath(dir, sha1 string) string {
	d1 := sha1[:2]
	d2 := sha1[2:4]
	return filepath.Join(dir, "blobs_crashes", d1, d2, sha1)
}

// MessageFilePath returns path of a message file
func (s *StoreCrashes) MessageFilePath(sha1 []byte) string {
	sha1Str := hex.EncodeToString(sha1)
	return blobCrahesPath(s.dataDir, sha1Str)
}

// MessageFileExists return true if message file exists
func (s *StoreCrashes) MessageFileExists(sha1 []byte) bool {
	p := s.MessageFilePath(sha1)
	return PathExists(p)
}

func (s *StoreCrashes) appendString(str string) error {
	_, err := s.dataFile.WriteString(str)
	if err != nil {
		logger.Errorf("StoreCrashes.appendString() error: %s\n", err)
	}
	return err
}

func (s *StoreCrashes) writeMessageAsSha1(msg []byte, sha1 []byte) error {
	path := s.MessageFilePath(sha1)
	err := WriteBytesToFile(msg, path)
	if err != nil {
		logger.Errorf("StoreCrashes.writeMessageAsSha1(): failed to write %s with error %s", path, err)
	}
	return err
}

func ip2str(s string) uint32 {
	var nums [4]uint32
	parts := strings.Split(s, ".")
	for n, p := range parts {
		num, _ := strconv.Atoi(p)
		nums[n] = uint32(num)
	}
	return (nums[0] << 24) | (nums[1] << 16) + (nums[2] << 8) | nums[3]
}

// GetApps returns all unique apps
func (s *StoreCrashes) GetApps() []*App {
	s.Lock()
	defer s.Unlock()
	return s.apps
}

// GetCrashesForApp returns crashes for an app
func (s *StoreCrashes) GetCrashesForApp(appName string) []*Crash {
	s.Lock()
	defer s.Unlock()
	app := s.FindOrCreateApp(appName)
	return app.Crashes
}

func (s *StoreCrashes) GetCrashesForIpAddrInternal(app *App, ipAddrInternal string) []*Crash {
	s.Lock()
	defer s.Unlock()
	res := make([]*Crash, 0)
	ipAddrPtr := s.FindIp(ipAddrInternal)
	if ipAddrPtr == nil {
		return res
	}
	for i, c := range s.crashes {
		if c.App == app && c.IpAddrInternal == ipAddrPtr {
			res = append(res, s.crashes[i])
		}
	}
	sort.Sort(CrashesByCreatedOn(res))
	return res
}

func (s *StoreCrashes) GetCrashesForCrashingLine(app *App, crashingLine string) []*Crash {
	s.Lock()
	defer s.Unlock()
	res := make([]*Crash, 0)
	crashingLinePtr := s.FindCrashingLine(crashingLine)
	if crashingLinePtr == nil {
		return res
	}
	for i, c := range s.crashes {
		if c.App == app && c.CrashingLine == crashingLinePtr {
			res = append(res, s.crashes[i])
		}
	}
	sort.Sort(CrashesByCreatedOn(res))
	return res
}

func (s *StoreCrashes) GetCrashById(id int) *Crash {
	s.Lock()
	defer s.Unlock()
	if id < 0 || id > len(s.crashes) {
		return nil
	}
	return s.crashes[id]
}

func serCrash(c *Crash) string {
	s1 := base64.StdEncoding.EncodeToString(c.Sha1[:])
	s1 = s1[:len(s1)-1] // remove '=' from the end
	s2 := fmt.Sprintf("%d", c.CreatedOn.Unix())
	s3 := remSep(c.App.Name)
	s4 := remSep(*c.ProgramVersion)
	s5 := *c.IpAddrInternal
	s6 := *c.CrashingLine
	return fmt.Sprintf("C%s|%s|%s|%s|%s|%s\n", s1, s2, s3, s4, s5, s6)
}

func getCrashPrefixData(crash *Crash) []byte {
	s := fmt.Sprintf("App: %s\n", crash.App.Name)
	s += fmt.Sprintf("Ip: %s\n", *crash.IpAddrInternal)
	s += fmt.Sprintf("On: %s\n", crash.CreatedOn.Format(time.RFC3339))
	return []byte(s)
}

func storeCrashingLine(d []byte) string {
	l := ExtractSumatraCrashingLine(d)
	return remSep(l)
}

func (s *StoreCrashes) SaveCrash(appName, appVer, ipAddr string, crashData []byte) error {
	s.Lock()
	defer s.Unlock()

	// TODO: white-list app names?
	app := s.FindOrCreateApp(appName)
	programVersionInterned := s.FindOrCreateVersion(appVer)
	ipAddrInterned := s.FindOrCreateIp(ipAddrToInternal(ipAddr))
	cl := storeCrashingLine(crashData)
	crashingLine := s.FindOrCreateCrashingLine(cl)

	c := &Crash{
		Id:             len(s.crashes),
		App:            app,
		CreatedOn:      time.Now(),
		ProgramVersion: programVersionInterned,
		IpAddrInternal: ipAddrInterned,
		CrashingLine:   crashingLine,
	}

	var buf bytes.Buffer
	buf.Write(getCrashPrefixData(c))
	buf.Write(crashData)
	dstData := buf.Bytes()
	sha1 := Sha1OfBytes(dstData)
	copy(c.Sha1[:], sha1)

	if err := s.writeMessageAsSha1(crashData, sha1); err != nil {
		return err
	}
	crashLine := serCrash(c)
	if err := s.appendString(crashLine); err != nil {
		return err
	}

	s.appendCrash(c)
	return nil
}
