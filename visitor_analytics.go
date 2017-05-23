package main

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kjk/dailyrotate"
	"github.com/kjk/siser"
)

const (
	keyURI      = "uri"
	keyCode     = "code"
	keyIPAddr   = "ip"
	keyWhen     = "when"
	keyDuration = "dur" // in milliseconds
	keyReferer  = "referer"
)

var (
	analyticsFile *dailyrotate.File
)

type countedString struct {
	s string
	n int
}

type analyticsStats struct {
	mostPopularURI     []countedString
	mostPopularReferer []countedString
	nUniqueIPs         int
}

func getMostCounted(m map[string]int, nTop int) []countedString {
	var res []countedString
	for s, n := range m {
		cs := countedString{
			s: s,
			n: n,
		}
		res = append(res, cs)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].n > res[j].n
	})
	if len(res) > nTop {
		return res[:nTop]
	}
	return res
}

// TODO: slowest pages
func calcAnalyticsStats(path string) *analyticsStats {
	uriCount := make(map[string]int)
	refererCount := make(map[string]int)
	ipCount := make(map[string]int)

	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	r := siser.NewReader(f)

	for r.ReadNext() {
		_, rec := r.Record()
		if requestURI, ok := rec.Get(keyURI); ok {
			uri, err := url.ParseRequestURI(requestURI)
			if err != nil {
				s := uri.Path
				uriCount[s]++
			}
		}
		if referer, ok := rec.Get(keyReferer); ok {
			// we don't care about internal referers
			if !strings.Contains(referer, "blog.kowalczyk.info") {
				refererCount[referer]++
			}
		}
		if ip, ok := rec.Get(keyIPAddr); ok {
			ipCount[ip]++
		}
	}
	if r.Err() != nil {
		return nil
	}
	return &analyticsStats{
		mostPopularURI:     getMostCounted(uriCount, 32),
		mostPopularReferer: getMostCounted(refererCount, 64),
		nUniqueIPs:         len(ipCount),
	}
}

func onAnalyticsFileCloseBackground(path string) {
	// TODO:
	// - compress
	// - upload to backblaze
	// - do basic stats
	// we don't want to block writes
}

func onAnalyticsFileClosed(path string, didRotate bool) {
	if didRotate {
		go onAnalyticsFileCloseBackground(path)
	}
}

func initAnalyticsMust(pathFormat string) error {
	var err error
	analyticsFile, err = dailyrotate.NewFile(pathFormat, onAnalyticsFileClosed)
	fatalIfErr(err)
	return nil
}

func logWebAnalytics(r *http.Request, code int, dur time.Duration) {
	uri := r.RequestURI

	// don't log hits we don't care about
	if uri == "/robots.txt" {
		return
	}
	ext := strings.ToLower(filepath.Ext(uri))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".ico", ".gif":
		return
	}

	ipAddr := getIPAddress(r)
	referer := r.Referer()
	when := time.Now().UTC().Format(time.RFC3339)
	codeStr := strconv.Itoa(code)
	durMs := float64(dur) / float64(time.Millisecond)
	durStr := strconv.FormatFloat(durMs, 'f', 2, 64)
	var rec siser.Record
	rec = rec.Append(keyURI, uri, keyCode, codeStr, keyIPAddr, ipAddr, keyReferer, referer, keyDuration, durStr, keyWhen, when)
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
