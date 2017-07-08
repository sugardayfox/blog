package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kjk/u"
)

func servePlainText(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

func serveXML(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "application/xml")
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

	s = fmt.Sprintf("Raw RemoteAddr: %s", r.RemoteAddr)
	a = append(a, s)

	s = fmt.Sprintf("Real RemoteAddr: %s", u.RequestGetRemoteAddress(r))
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
	a = append(a, fmt.Sprintf("ver: https://github.com/kjk/blog/commit/%s", sha1ver))

	s = strings.Join(a, "\n")
	servePlainText(w, s)
}

// GET /app/sendmsg
// args:
//  msg: text to email me
//  cookie: must match randomCookie
func handleSendMsg(w http.ResponseWriter, r *http.Request) {
	msg := strings.TrimSpace(r.FormValue("msg"))
	cookie := strings.TrimSpace(r.FormValue("cookie"))
	if cookie == "" {
		logger.Notice("handleSendMsg: 'cookie' arg is missing\n")
		// technically that should be some other error
		serve404(w, r)
		return
	}
	if cookie != randomCookie {
		logger.Noticef("handleSendMsg: 'cookie' != randomCookie (%s != %s)\n", cookie, randomCookie)
		// technically that should be some other error
		serve404(w, r)
		return
	}
	if msg == "" {
		logger.Notice("handleSendMsg: 'msg' arg is missing\n")
		// technically that should be some other error
		serve404(w, r)
		return
	}
	// I assume that shorter messages are garbage
	words := strings.Split(msg, " ")
	if len(words) < 3 {
		logger.Noticef("handleSendMsg: 'msg' is too short ('%s')\n", msg)
		// technically that should be some other error
		serve404(w, r)
		return
	}
	sendMail("Message from contact me blog page", msg)
	logger.Noticef("handleSendMsg: sent email with message: '%s'\n", msg)
	servePlainText(w, "ok")
}
