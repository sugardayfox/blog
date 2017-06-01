package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/kjk/dailyrotate"
	"github.com/kjk/siser"
)

const (
	keyURI       = "uri"
	keyCode      = "code"
	keyIPAddr    = "ip"
	keyWhen      = "when"
	keyDuration  = "dur" // in milliseconds
	keyReferer   = "referer"
	keySize      = "size"
	keyUserAgent = "ua"
)

var (
	analyticsFile *dailyrotate.File
)

type countedString struct {
	s string
	n int
}

type analyticsStats struct {
	urls       []countedString
	referers   []countedString
	notFound   []countedString
	nUniqueIPs int
}

func initAnalyticsMust(pathFormat string) error {
	var err error
	analyticsFile, err = dailyrotate.NewFile(pathFormat, onAnalyticsFileClosed)
	fatalIfErr(err)
	return nil
}

func withAnalyticsLogging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		rrw := NewRecordingResponseWriter(w)
		f(rrw, r)
		dur := time.Since(timeStart)
		logWebAnalytics(r, rrw.Code, rrw.BytesWritten, dur)
	}
}

func countedStringMapToArray(m map[string]int) []countedString {
	var res []countedString
	for s, n := range m {
		cs := countedString{
			s: s,
			n: n,
		}
		res = append(res, cs)
	}
	// sort in reverse: most frequent first
	sort.Slice(res, func(i, j int) bool {
		return res[i].n > res[j].n
	})
	return res
}

// TODO:
// - slowest pages
func calcAnalyticsStats(path string) (*analyticsStats, error) {
	uriCount := make(map[string]int)
	uri404Count := make(map[string]int)
	refererCount := make(map[string]int)
	ipCount := make(map[string]int)

	f, err := openFileMaybeCompressed(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := siser.NewReader(f)

	for r.ReadNext() {
		_, rec := r.Record()
		code, ok1 := rec.Get(keyCode)
		requestURI, ok2 := rec.Get(keyURI)
		ip, ok3 := rec.Get(keyIPAddr)
		// referer can be empty
		referer, _ := rec.Get(keyReferer)

		if !(ok1 && ok2 && ok3) {
			// shouldn't happen
			continue
		}
		uri, err := url.ParseRequestURI(requestURI)
		if err != nil {
			// shouldn't happen
			continue
		}

		if code == "404" {
			uri404Count[uri.Path]++
			continue
		}

		// we don't care about internal referers
		if referer != "" && !strings.Contains(referer, "blog.kowalczyk.info") {
			refererCount[referer]++
		}

		// don't record redirects
		if code != "200" {
			continue
		}
		uriCount[uri.Path]++
		ipCount[ip]++
	}
	if r.Err() != nil {
		return nil, r.Err()
	}
	return &analyticsStats{
		urls:       countedStringMapToArray(uriCount),
		referers:   countedStringMapToArray(refererCount),
		notFound:   countedStringMapToArray(uri404Count),
		nUniqueIPs: len(ipCount),
	}, nil
}

func analyticsStatsText(a *analyticsStats) []string {
	if a == nil {
		return []string{"Couldn't calculate analytics stats"}
	}

	var lines []string
	s := fmt.Sprintf("Unique ips: %d, unique referers: %d, unique urls: %d", a.nUniqueIPs, len(a.referers), len(a.urls))
	lines = append(lines, s)

	lines = append(lines, "\nMost frequent referers:")
	n := len(a.referers)
	if n > 64 {
		n = 64
	}
	for i := 0; i < n; i++ {
		cs := a.referers[i]
		s = fmt.Sprintf("%s : %d", cs.s, cs.n)
		lines = append(lines, s)
	}

	lines = append(lines, "\nMost popular urls:")
	n = len(a.urls)
	if n > 64 {
		n = 64
	}
	for i := 0; i < n; i++ {
		cs := a.urls[i]
		s = fmt.Sprintf("%s : %d", cs.s, cs.n)
		lines = append(lines, s)
	}
	return lines
}

func onAnalyticsFileCloseBackground(path string) {
	timeStart := time.Now()
	a, statsErr := calcAnalyticsStats(path)
	dur := time.Since(timeStart)
	var lines []string
	size, _ := getFileSize(path)
	sizeStr := humanize.Bytes(uint64(size))

	timeStart = time.Now()
	dstPath := path + ".gz"
	err := gzipFile(dstPath, path)
	if err != nil {
		s := fmt.Sprintf("gzipFile(%s, %s) failed with %s", dstPath, path, err)
		lines = append(lines, s)
	}
	durCompress := time.Since(timeStart)
	os.Remove(path)

	fileName := filepath.Base(dstPath)
	b2Path := "backup/blog/analytics/" + fileName
	timeStart = time.Now()
	b2UploadFile(b2Path, dstPath)
	durUpload := time.Since(timeStart)

	s := fmt.Sprintf("Processing analytics for %s of size %s took %s. Compressing took %s. Uploading to b2 as %s took %s.", path, sizeStr, dur, durCompress, b2Path, durUpload)
	lines = append(lines, s)
	if statsErr != nil {
		s = fmt.Sprintf("Processing analytics failed with %s", statsErr)
		lines = append(lines, s)
	} else {
		lines = append(lines, analyticsStatsText(a)...)
	}
	subject := utcNow().Format("blog stats on 2006-01-02 15:04:05")
	body := strings.Join(lines, "\n\n")
	sendMail(subject, body)
}

func onAnalyticsFileClosed(path string, didRotate bool) {
	logger.Noticef("onAnalyticsFileClosed: %s, didRotate: %v\n", path, didRotate)
	if didRotate {
		// do in background, we don't want to block writes
		go onAnalyticsFileCloseBackground(path)
	}
}

func logWebAnalytics(r *http.Request, code int, nBytesWritten int64, dur time.Duration) {
	uri := r.RequestURI

	// don't log hits we don't care about
	if uri == "/robots.txt" {
		return
	}
	ext := strings.ToLower(filepath.Ext(uri))
	switch ext {
	// we care mostly about .http files, those are referenced files
	case ".png", ".jpg", ".jpeg", ".ico", ".gif", ".css", ".js":
		return
	}

	ipAddr := getIPAddress(r)
	when := time.Now().UTC().Format(time.RFC3339)
	codeStr := strconv.Itoa(code)
	durMs := float64(dur) / float64(time.Millisecond)
	durStr := strconv.FormatFloat(durMs, 'f', 2, 64)
	sizeStr := strconv.FormatInt(nBytesWritten, 10)
	var rec siser.Record
	rec = rec.Append(keyURI, uri, keyCode, codeStr, keyIPAddr, ipAddr, keyDuration, durStr, keyWhen, when, keySize, sizeStr)

	ua := r.Header.Get("user-agent")
	if ua != "" {
		rec = rec.Append(keyUserAgent, ua)
	}
	referer := r.Referer()
	if referer != "" {
		rec = rec.Append(keyReferer, referer)
	}
	d := rec.Marshal()
	// ignoring error because can't do anything about it
	analyticsFile.Write2(d, true)
}

func analyticsClose() {
	if analyticsFile != nil {
		analyticsFile.Close()
		analyticsFile = nil
	}
}

func testAnalyticsStats(path string) {
	stats, err := calcAnalyticsStats(path)
	fatalIfErr(err)
	lines := analyticsStatsText(stats)
	fmt.Printf("Analytics as text:\n%s\n", strings.Join(lines, "\n"))
}
