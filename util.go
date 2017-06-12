// This code is under BSD license. See license-bsd.txt
package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	patWs            = regexp.MustCompile(`\s+`)
	patCharsToRemove = regexp.MustCompile("[+:*%&/()]")
)

// urlify generates safe url from tile
func urlify(title string) string {
	s := strings.TrimSpace(title)
	s = strings.ToLower(s)
	s = patCharsToRemove.ReplaceAllString(s, "")
	s = patWs.ReplaceAllString(s, "-")
	if len(s) > 128 {
		s = s[:128]
	}
	return s
}

func httpErrorf(w http.ResponseWriter, format string, args ...interface{}) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	http.Error(w, msg, http.StatusBadRequest)
}
