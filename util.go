// This code is under BSD license. See license-bsd.txt
package main

import (
	"crypto/sha1"
	"fmt"
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

const base64Chars = "0123456789abcdefghijklmnopqrstuvwxyz"

// ShortenID encodes n as base64
func ShortenID(n int) string {
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

// pathExists returns true if a filesystem path exists
// Treats any error (e.g. lack of access due to permissions) as non-existence
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getFileSize(path string) (int64, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
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

func utcNow() time.Time {
	return time.Now().UTC()
}
