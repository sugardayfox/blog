// This code is under BSD license. See license-bsd.txt
package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
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

// PanicIfErr panics if err is not nil
func PanicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// PanicIf panics if cond is true
func PanicIf(cond bool, args ...interface{}) {
	if !cond {
		return
	}
	msg := "invalid state"
	if len(args) > 0 {
		s, ok := args[0].(string)
		if ok {
			msg = s
			if len(s) > 1 {
				msg = fmt.Sprintf(msg, args[1:]...)
			}
		}
	}
	panic(msg)
}

func panicif(cond bool, args ...interface{}) {
	if !cond {
		return
	}
	msg := "panic"
	if len(args) > 0 {
		s, ok := args[0].(string)
		if ok {
			msg = s
			if len(s) > 1 {
				msg = fmt.Sprintf(msg, args[1:]...)
			}
		}
	}
	panic(msg)
}

// Urlify generates url from tile
func Urlify(title string) string {
	s := strings.TrimSpace(title)
	s = patWs.ReplaceAllString(s, "-")
	s = patNonAlpha.ReplaceAllString(s, "")
	s = patMultipleMinus.ReplaceAllString(s, "-")
	if len(s) > 48 {
		s = s[:48]
	}
	return s
}

const (
	cr = 0xd
	lf = 0xa
)

// ExtracLine finds end of line (cr, lf or crlf). Return the line
// and the remaining of data (without the end-of-line character(s))
func ExtractLine(d []byte) ([]byte, []byte) {
	if d == nil || len(d) == 0 {
		return nil, nil
	}
	wasCr := false
	pos := -1
	for i := 0; i < len(d); i++ {
		if d[i] == cr || d[i] == lf {
			wasCr = (d[i] == cr)
			pos = i
			break
		}
	}
	if pos == -1 {
		return d, nil
	}
	line := d[:pos]
	rest := d[pos+1:]
	if wasCr && len(rest) > 0 && rest[0] == lf {
		rest = rest[1:]
	}
	return line, rest
}

// SkipPastLine iterates d as lines, find lineToFind and return the part
// after that line. Return nil if not found
func SkipPastLine(d []byte, lineToFind string) []byte {
	lb := []byte(lineToFind)
	var l []byte
	for {
		l, d = ExtractLine(d)
		if l == nil {
			return nil
		}
		if bytes.Equal(l, lb) {
			return d
		}
	}
}

// FindLineWithPrefix finds a line with a given prefix
func FindLineWithPrefix(d []byte, prefix string) []byte {
	prefixb := []byte(prefix)
	var l []byte
	for {
		l, d = ExtractLine(d)
		if l == nil {
			return nil
		}
		if bytes.HasPrefix(l, prefixb) {
			return l
		}
	}
}

const base64Chars = "0123456789abcdefghijklmnopqrstuvwxyz"

// ShortenId encodes n as base64
func ShortenId(n int) string {
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

// UnshortenID decodes base64 string
func UnshortenID(s string) int {
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

// DurationToString converts duration to a string
func DurationToString(d time.Duration) string {
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

// TimeSinceNowAsString returns string version of time since a ginve timestamp
func TimeSinceNowAsString(t time.Time) string {
	return DurationToString(time.Now().Sub(t))
}

// Sha1HexOfBytes returns 40-byte hex sha1 of bytes
func Sha1HexOfBytes(data []byte) string {
	return fmt.Sprintf("%x", Sha1OfBytes(data))
}

// Sha1OfBytes returns 20-byte sha1 of bytes
func Sha1OfBytes(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// PathExists returns true if a filesystem path exists
// Treats any error (e.g. lack of access due to permissions) as non-existence
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ListFilesInDir returns a list of files in a directory
func ListFilesInDir(dir string, recursive bool) []string {
	files := make([]string, 0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		isDir, err := PathIsDir(path)
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

// UserHomeDir returns $HOME diretory of the user
func UserHomeDir() string {
	// user.Current() returns nil if cross-compiled e.g. on mac for linux
	if usr, _ := user.Current(); usr != nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

// ExpandTildeInPath converts ~ to $HOME
func ExpandTildeInPath(s string) string {
	if strings.HasPrefix(s, "~") {
		return UserHomeDir() + s[1:]
	}
	return s
}

// CreateDirForFileMust is like CreateDirForFile. Panics on error.
func CreateDirForFileMust(path string) string {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	PanicIfErr(err)
	return dir
}

// PathIsDir returns true if a path exists and is a directory
// Returns false, nil if a path exists and is not a directory (e.g. a file)
// Returns undefined, error if there was an error e.g. because a path doesn't exists
func PathIsDir(path string) (isDir bool, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// WriteBytesToFile is like ioutil.WriteFile() but also creates intermediary
// directories
func WriteBytesToFile(d []byte, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, d, 0644)
}
