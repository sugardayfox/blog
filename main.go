package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/acme/autocert"
)

var (
	analyticsCode = "UA-194516-1"

	logger  *ServerLogger
	dataDir string
	store   *ArticlesStore
	sha1ver string
)

func getDataDir() string {
	if dataDir != "" {
		return dataDir
	}

	dirsToCheck := []string{"/data", expandTildeInPath("~/data/blog")}
	for _, dir := range dirsToCheck {
		if pathExists(dir) {
			dataDir = dir
			return dataDir
		}
	}

	log.Fatalf("data directory (%v) doesn't exist", dirsToCheck)
	return ""
}

func isTopLevelURL(url string) bool {
	return 0 == len(url) || "/" == url
}

func getReferer(r *http.Request) string {
	return r.Header.Get("Referer")
}

// this list was determined by watching /logs
var noLog404 = map[string]bool{
	"/crossdomain.xml":                                               true,
	"/article/Exercise-links-1.html":                                 true,
	"/article/Ecco-for-free.html":                                    true,
	"/article/Disappointed-by-The-Bat.html":                          true,
	"/article/Comments-need-not-apply.html":                          true,
	"/article/Browsing-Newton.html":                                  true,
	"/article/Perl-and-lisp-programmers.html":                        true,
	"/article/iPod-competition.html":                                 true,
	"/article/Programming-Jabber.html":                               true,
	"/article/Good-software-design-contradicts-eXtreme-Program.html": true,
	"/article/Bloglines-vs-Google-Reader-the-verdict.html":           true,
	"/2002/07/30/stuid-coding-mistake-of-the-day.html":               true,
	"/article/Corman-Lisp.html":                                      true,
	"/article/Offshore-outsourcing.html":                             true,
	"/article/Nabble-hosted-forums.html":                             true,
}

func shouldLog404(s string) bool {
	if strings.HasPrefix(s, "/apple-touch-icon") {
		return false
	}
	_, ok := noLog404[s]
	return !ok
}

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func getIPAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

func setContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

func writeResponse(w http.ResponseWriter, responseBody string) {
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(responseBody)), 10))
	io.WriteString(w, responseBody)
}

func textResponse(w http.ResponseWriter, text string) {
	setContentType(w, "text/plain")
	writeResponse(w, text)
}

var (
	flgHTTPAddr        string
	flgProduction      bool
	flgNewArticleTitle string
)

func parseCmdLineFlags() {
	flag.StringVar(&flgHTTPAddr, "addr", ":5020", "HTTP server address")
	flag.BoolVar(&flgProduction, "production", false, "are we running in production")
	flag.StringVar(&flgNewArticleTitle, "newarticle", "", "create a new article")
	flag.Parse()
}

func isTmpFile(path string) bool {
	return strings.HasSuffix(path, ".tmp")
}

func sanitizeForFile(s string) string {
	var res []byte
	toRemove := "/\\#()[]{},?+.'\""
	var prev rune
	buf := make([]byte, 3)
	for _, c := range s {
		if strings.ContainsRune(toRemove, c) {
			continue
		}
		switch c {
		case ' ', '_':
			c = '-'
		}
		if c == prev {
			continue
		}
		prev = c
		n := utf8.EncodeRune(buf, c)
		for i := 0; i < n; i++ {
			res = append(res, buf[i])
		}
	}
	if len(res) > 32 {
		res = res[:32]
	}
	s = string(res)
	s = strings.Trim(s, "_- ")
	s = strings.ToLower(s)
	return s
}

func findUniqueArticleID(articles []*Article) int {
	var ids []int
	for _, a := range articles {
		ids = append(ids, a.ID)
	}
	if len(ids) == 0 {
		return 1
	}
	sort.Ints(ids)
	prevID := ids[0]
	for i := 1; i < len(ids); i++ {
		if ids[i] != prevID+1 {
			return prevID + 1
		}
		prevID = ids[i]
	}
	return prevID + 1
}

func genNewArticle(title string) {
	fmt.Printf("genNewArticle: %q\n", title)
	store, err := NewArticlesStore()
	if err != nil {
		log.Fatalf("NewStore() failed with %s", err)
	}
	newID := findUniqueArticleID(store.articles)
	t := time.Now()
	dir := "blog_posts"
	yyyy := fmt.Sprintf("%04d", t.Year())
	month := t.Month()
	sanitizedTitle := sanitizeForFile(title)
	name := fmt.Sprintf("%02d-%s.md", month, sanitizedTitle)
	fmt.Printf("new id: %d, name: %s\n", newID, name)
	path := filepath.Join(dir, yyyy, name)
	s := fmt.Sprintf(`Id: %d
Title: %s
Date: %s
Format: Markdown
--------------`, newID, title, t.Format(time.RFC3339))
	for i := 1; i < 10; i++ {
		if !pathExists(path) {
			break
		}
		name = fmt.Sprintf("%02d-%s-%d.md", month, sanitizedTitle, i)
		path = filepath.Join(dir, yyyy, name)
	}
	fatalIf(pathExists(path))
	fmt.Printf("path: %s\n", path)
	createDirForFileMust(path)
	ioutil.WriteFile(path, []byte(s), 0644)
}

func loadArticles() {
	var err error
	if store, err = NewArticlesStore(); err != nil {
		log.Fatalf("NewStore() failed with %s", err)
	}
	articles := store.GetArticles()
	articlesCache.articles = articles
	articlesCache.articlesJs, articlesCache.articlesJsSha1 = buildArticlesJSON(articles)
}

func hostPolicy(ctx context.Context, host string) error {
	if strings.HasSuffix(host, "kowalczyk.info") {
		return nil
	}
	return errors.New("acme/autocert: only *.kowalczyk.info hosts are allowed")
}

func main() {
	parseCmdLineFlags()

	if false {
		testAnalyticsStats("/Users/kjk/Downloads/2017-06-02.txt.gz")
		os.Exit(0)
	}

	if flgNewArticleTitle != "" {
		genNewArticle(flgNewArticleTitle)
		return
	}

	if flgProduction {
		reloadTemplates = false
		flgHTTPAddr = ":80"
	} else {
		analyticsCode = ""
	}

	logger = NewServerLogger(256, 256)

	rand.Seed(time.Now().UnixNano())

	loadArticles()

	readRedirects()

	analyticsPath := filepath.Join(getDataDir(), "analytics", "2006-01-02.txt")
	initAnalyticsMust(analyticsPath)

	ctx := context.TODO()
	var wg sync.WaitGroup
	var httpsSrv, httpSrv *http.Server

	if flgProduction {
		httpsSrv = makeHTTPServer()
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
		}
		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		logger.Noticef("Starting https server on %s\n", httpsSrv.Addr)
		go func() {
			wg.Add(1)
			err := httpsSrv.ListenAndServeTLS("", "")
			// mute error caused by Shutdown()
			if err == http.ErrServerClosed {
				err = nil
			}
			fatalIfErr(err)
			fmt.Printf("HTTPS server shutdown gracefully\n")
			wg.Done()
		}()
	}

	httpSrv = makeHTTPServer()
	httpSrv.Addr = flgHTTPAddr
	logger.Noticef("Starting http server on %s, in production: %v, sha1ver: %s", httpSrv.Addr, flgProduction, sha1ver)
	go func() {
		wg.Add(1)
		err := httpSrv.ListenAndServe()
		// mute error caused by Shutdown()
		if err == http.ErrServerClosed {
			err = nil
		}
		fatalIfErr(err)
		fmt.Printf("HTTP server shutdown gracefully\n")
		wg.Done()
	}()

	if flgProduction {
		sendBootMail()
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt /* SIGINT */, syscall.SIGTERM)
	sig := <-c
	fmt.Printf("Got signal %s\n", sig)
	if httpsSrv != nil {
		httpsSrv.Shutdown(ctx)
	}
	if httpSrv != nil {
		httpSrv.Shutdown(ctx)
	}
	wg.Wait()
	analyticsClose()
	fmt.Printf("Exited\n")
}
