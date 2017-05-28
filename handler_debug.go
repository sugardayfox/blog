package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func servePlainText(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

// /app/debug
func handleDebug(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("url: %s %s", r.Method, r.RequestURI)
	a := []string{s}

	s = "https: no"
	if r.TLS != nil {
		s = "https: yes"
	}
	a = append(a, s)

	s = fmt.Sprintf("RemoteAddr: %s", r.RemoteAddr)
	a = append(a, s)

	a = append(a, "Headers:")
	for k, v := range r.Header {
		if len(v) == 0 {
			a = append(a, k)
		} else if len(v) == 1 {
			s = fmt.Sprintf("  %s: %v", k, v[0])
			a = append(a, s)
		} else {
			a = append(a, "  "+k+":")
			for _, v2 := range v {
				a = append(a, "    "+v2)
			}
		}
	}

	a = append(a, "")
	a = append(a, fmt.Sprintf("ver: https://github.com/kjk/web-blog/commit/%s", sha1ver))

	s = strings.Join(a, "\n")
	servePlainText(w, s)
}
