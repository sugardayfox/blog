// This code is under BSD license. See license-bsd.txt
package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var patWs = regexp.MustCompile(`\s+`)
var patNonAlpha = regexp.MustCompile(`[^\w-]`)
var patMultipleMinus = regexp.MustCompile("-+")

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

func httpErrorf(w http.ResponseWriter, format string, args ...interface{}) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	http.Error(w, msg, http.StatusBadRequest)
}
