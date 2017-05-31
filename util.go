// This code is under BSD license. See license-bsd.txt
package main

import (
	"compress/bzip2"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var patWs = regexp.MustCompile(`\s+`)
var patNonAlpha = regexp.MustCompile(`[^\w-]`)
var patMultipleMinus = regexp.MustCompile("-+")

func fmtArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	format := args[0].(string)
	if len(args) == 1 {
		return format
	}
	return fmt.Sprintf(format, args[1:]...)
}

func panicWithMsg(defaultMsg string, args ...interface{}) {
	s := fmtArgs(args...)
	if s == "" {
		s = defaultMsg
	}
	fmt.Printf("%s\n", s)
	panic(s)
}

func fatalIfErr(err error, args ...interface{}) {
	if err == nil {
		return
	}
	panicWithMsg(err.Error(), args...)
}

func fatalIf(cond bool, args ...interface{}) {
	if !cond {
		return
	}
	panicWithMsg("fatalIf: condition failed", args...)
}

// urlify generates url from tile
func urlify(title string) string {
	s := strings.TrimSpace(title)
	s = patWs.ReplaceAllString(s, "-")
	s = patNonAlpha.ReplaceAllString(s, "")
	s = patMultipleMinus.ReplaceAllString(s, "-")
	s = strings.Replace(s, ":", "", -1)
	s = strings.Replace(s, "%", "-perc", -1)
	if len(s) > 48 {
		s = s[:48]
	}
	return s
}

const base64Chars = "0123456789abcdefghijklmnopqrstuvwxyz"

// shortenID encodes n as base64
func shortenID(n int) string {
	var buf [16]byte
	size := 0
	for {
		buf[size] = base64Chars[n%36]
		size++
		if n < 36 {
			break
		}
		n /= 36
	}
	end := size - 1
	for i := 0; i < end; i++ {
		b := buf[i]
		buf[i] = buf[end]
		buf[end] = b
		end--
	}
	return string(buf[:size])
}

// unshortenID decodes base64 string
func unshortenID(s string) int {
	n := 0
	for _, c := range s {
		n *= 36
		i := strings.IndexRune(base64Chars, c)
		// TODO: return an error if i == -1
		n += i
	}
	return n
}

func httpErrorf(w http.ResponseWriter, format string, args ...interface{}) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	http.Error(w, msg, http.StatusBadRequest)
}

// durationToString converts duration to a string
func durationToString(d time.Duration) string {
	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours())
	days := hours / 24
	hours = hours % 24
	if days > 0 {
		return fmt.Sprintf("%dd %dhr", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dhr %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func timeSinceNowAsString(t time.Time) string {
	return durationToString(time.Now().Sub(t))
}

// sha1HexOfBytes returns 40-byte hex sha1 of bytes
func sha1HexOfBytes(data []byte) string {
	return fmt.Sprintf("%x", sha1OfBytes(data))
}

// sha1OfBytes returns 20-byte sha1 of bytes
func sha1OfBytes(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// pathExists returns true if a filesystem path exists
// Treats any error (e.g. lack of access due to permissions) as non-existence
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// pathIsDir returns true if a path exists and is a directory
// Returns false, nil if a path exists and is not a directory (e.g. a file)
// Returns undefined, error if there was an error e.g. because a path doesn't exists
func pathIsDir(path string) (isDir bool, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func getFileSize(path string) (int64, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func gzipFile(dstPath, srcPath string) error {
	fSrc, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fSrc.Close()
	fDst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer fDst.Close()
	w, err := gzip.NewWriterLevel(fDst, gzip.BestCompression)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, fSrc)
	if err != nil {
		return err
	}
	return nil
}

func listFilesInDir(dir string, recursive bool) []string {
	files := make([]string, 0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		isDir, err := pathIsDir(path)
		if err != nil {
			return err
		}
		if isDir {
			if recursive || path == dir {
				return nil
			}
			return filepath.SkipDir
		}
		files = append(files, path)
		return nil
	})
	return files
}

// userHomeDir returns $HOME diretory of the user
func userHomeDir() string {
	// user.Current() returns nil if cross-compiled e.g. on mac for linux
	if usr, _ := user.Current(); usr != nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

// expandTildeInPath converts ~ to $HOME
func expandTildeInPath(s string) string {
	if strings.HasPrefix(s, "~") {
		return userHomeDir() + s[1:]
	}
	return s
}

func createDirForFileMust(path string) string {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	fatalIfErr(err)
	return dir
}

// implement io.ReadCloser over os.File wrapped with io.Reader.
// io.Closer goes to os.File, io.Reader goes to wrapping reader
type readerWrappedFile struct {
	f *os.File
	r io.Reader
}

func (rc *readerWrappedFile) Close() error {
	return rc.f.Close()
}

func (rc *readerWrappedFile) Read(p []byte) (int, error) {
	return rc.r.Read(p)
}

func openFileMaybeCompressed(path string) (io.ReadCloser, error) {
	ext := strings.ToLower(filepath.Ext(path))
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if ext == ".gz" {
		r, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		rc := &readerWrappedFile{
			f: f,
			r: r,
		}
		return rc, nil
	}
	if ext == ".bz2" {
		r := bzip2.NewReader(f)
		rc := &readerWrappedFile{
			f: f,
			r: r,
		}
		return rc, nil
	}
	return f, nil
}

func utcNow() time.Time {
	return time.Now().UTC()
}
