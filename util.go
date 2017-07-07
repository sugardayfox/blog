// This code is under BSD license. See license-bsd.txt
package main

import (
	"fmt"
	"net/http"
	"strings"
)

// whitelisted characters valid in url
func validateRune(c rune) byte {
	if c >= 'a' && c <= 'z' {
		return byte(c)
	}
	if c >= '0' && c <= '9' {
		return byte(c)
	}
	if c == '-' || c == '_' || c == '.' {
		return byte(c)
	}
	if c == ' ' {
		return '-'
	}
	return 0
}

func charCanRepeat(c byte) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

// urlify generates safe url from tile by removing hazardous characters
func urlify(title string) string {
	s := strings.TrimSpace(title)
	s = strings.ToLower(s)
	var res []byte
	for _, r := range s {
		c := validateRune(r)
		if c == 0 {
			continue
		}
		// eliminute duplicate consequitive characters
		var prev byte
		if len(res) > 0 {
			prev = res[len(res)-1]
		}
		if c == prev && !charCanRepeat(c) {
			continue
		}
		res = append(res, c)
	}
	s = string(res)
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

// RequestGetProtocol returns protocol under which the request is being served i.e. "http" or "https"
func RequestGetProtocol(r *http.Request) string {
	hdr := r.Header
	// X-Forwarded-Proto is set by proxies e.g. CloudFlare
	forwardedProto := strings.TrimSpace(strings.ToLower(hdr.Get("X-Forwarded-Proto")))
	if forwardedProto != "" {
		if forwardedProto == "http" || forwardedProto == "https" {
			return forwardedProto
		}
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

// RequestGetFullHost returns full host name e.g. "https://blog.kowalczyk.info/"
func RequestGetFullHost(r *http.Request) string {
	return RequestGetProtocol(r) + "://" + r.Host
}
