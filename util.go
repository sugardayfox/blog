// This code is under BSD license. See license-bsd.txt
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var patWs = regexp.MustCompile(`\s+`)
var patNonAlpha = regexp.MustCompile(`[^\w-]`)
var patMultipleMinus = regexp.MustCompile("-+")

// given an article title, generate a url
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

// find a end of line (cr, lf or crlf). Return the line
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

// iterate d as lines, find lineToFind and return the part
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
		end -= 1
	}
	return string(buf[:size])
}

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
