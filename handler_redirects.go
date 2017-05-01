package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var redirects = map[string]string{
	"/index.html":                  "/",
	"/blog":                        "/",
	"/blog/":                       "/",
	"/kb/serialization-in-c#.html": "/article/Serialization-in-C.html",
	"/extremeoptimizations":        "/extremeoptimizations/index.html",
	"/extremeoptimizations/":       "/extremeoptimizations/index.html",
	"/feed/rss2/atom.xml":          "/atom.xml",
	"/feed/rss2/":                  "/atom.xml",
	"/feed/rss2":                   "/atom.xml",
	"/feed/":                       "/atom.xml",
	"/feed":                        "/atom.xml",
	"/articles/cocoa-objectivec-reference.html":     "/articles/cocoa-reference.html",
	"/forum_sumatra/rss.php":                        "http://forums.fofou.org/sumatrapdf/rss",
	"/forum_sumatra":                                "http://forums.fofou.org/sumatrapdf",
	"/google6dba371684d43cd6.html":                  "/static/google6dba371684d43cd6.html",
	"/software/15minutes/index.html":                "/software/15minutes.html",
	"/software/15minutes/":                          "/software/15minutes.html",
	"/software/fofou":                               "/software/fofou/index.html",
	"/software/patheditor":                          "/software/patheditor/for-windows.html",
	"/software/patheditor/":                         "/software/patheditor/for-windows.html",
	"/software/scdiff/":                             "/software/scdiff.html",
	"/software/scdiff/index.html":                   "/software/scdiff.html",
	"/software/sumatra":                             "https://www.sumatrapdfreader.org/free-pdf-reader.html",
	"/software/sumatrapdf":                          "https://www.sumatrapdfreader.org/free-pdf-reader.html",
	"/software/sumatrapdf/":                         "https://www.sumatrapdfreader.org/free-pdf-reader.html",
	"/software/sumatrapdf/index.html":               "https://www.sumatrapdfreader.org/free-pdf-reader.html",
	"/software/sumatrapdf/download.html _blank":     "https://www.sumatrapdfreader.org/free-pdf-reader.html",
	"/software/sumatrapdf/download.html":            "https://www.sumatrapdfreader.org/free-pdf-reader.html",
	"/software/sumatrapdf/prerelase.html":           "/software/sumatrapdf/prerelease.html",
	"/software/sumatrapdf/sumatra-shot-00.gif":      "http://kjkpub.s3.amazonaws.com/blog/sumatra/sumatra-shot-00.gif",
	"/software/sumatrapdf/sumatra-shot-01.gif":      "http://kjkpub.s3.amazonaws.com/blog/sumatra/sumatra-shot-01.gif",
	"/software/sumatrapdf/sumatra-shot-00-full.gif": "http://kjkpub.s3.amazonaws.com/blog/sumatra/sumatra-shot-00-full.gif",
	"/software/sumatrapdf/sumatra-shot-01-full.gif": "http://kjkpub.s3.amazonaws.com/blog/sumatra/sumatra-shot-01-full.gif",
	"/software/sumatrapdf/SumatraSplash.png":        "http://kjkpub.s3.amazonaws.com/blog/sumatra/SumatraSplash.png",
	"/software/volante":                             "/software/volante/database.html",
	"/software/volante/":                            "/software/volante/database.html",
	"/software/volante/index.html":                  "/software/volante/database.html",
	"/software/fotofi":                              "/software/fotofi/free-stock-photos.html",
	"/software/fotofi/":                             "/software/fotofi/free-stock-photos.html",
	"/software/fotofi/index.html":                   "/software/fotofi/free-stock-photos.html",
	"/static/software.html":                         "/software/index.html",
	"/static/krzysztof.html":                        "/static/resume.html",
}

var articleRedirects = make(map[string]int)
var articleRedirectsMutex sync.Mutex

func readRedirects() {
	fname := filepath.Join("article_redirects.txt")
	d, err := ioutil.ReadFile(fname)
	if err != nil {
		return
	}
	lines := bytes.Split(d, []byte{'\n'})
	for _, l := range lines {
		if 0 == len(l) {
			continue
		}
		parts := strings.Split(string(l), "|")
		fatalIf(len(parts) != 2, "malformed article_redirects.txt, len(parts) = %d (!2)", len(parts))
		idStr := parts[0]
		url := strings.TrimSpace(parts[1])
		if id, err := strconv.Atoi(idStr); err != nil {
			panic("malformed line in article_redirects.txt")
		} else {
			a := store.GetArticleByID(id)
			if a != nil {
				articleRedirects[url] = id
				continue
			}
			//fmt.Printf("skipping redirect '%s' because article with id %d no longer present\n", string(l), id)
		}
	}
	logger.Noticef("loaded %d article redirects", len(articleRedirects))
}

// return -1 if there's no redirect for this urls
func getRedirectArticleID(url string) int {
	url = url[1:] // remove '/' from the beginning
	articleRedirectsMutex.Lock()
	defer articleRedirectsMutex.Unlock()
	if articleID, ok := articleRedirects[url]; ok {
		return articleID
	}
	return -1
}

func redirectIfNeeded(w http.ResponseWriter, r *http.Request) bool {
	uri := r.URL.Path
	//logger.Noticef("redirectIfNeeded(): %q", uri)

	if strings.HasPrefix(uri, "/software/sumatrapdf") {
		redirURL := "https://www.sumatrapdfreader.org" + uri[len("/software/sumatrapdf"):]
		http.Redirect(w, r, redirURL, 302)
		return true
	}

	if redirURL, ok := redirects[uri]; ok {
		//logger.Noticef("Redirecting %q => %q", url, redirUrl)
		http.Redirect(w, r, redirURL, 302)
		return true
	}

	redirectArticleID := getRedirectArticleID(uri)
	if redirectArticleID == -1 {
		return false
	}
	article := store.GetArticleByID(redirectArticleID)
	if article != nil {
		redirURL := "/" + article.Permalink()
		//logger.Noticef("Redirecting %q => %q", url, redirUrl)
		http.Redirect(w, r, redirURL, 302)
		return true
	}

	return false
}

// url: /forum_sumatra/${rest}
func forumRedirect(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/forum_sumatra/"):]

	redirURL := "http://forums.fofou.org/sumatrapdf/" + url
	if len(r.URL.RawQuery) > 0 {
		redirURL = redirURL + "?" + r.URL.RawQuery
	}
	//logger.Noticef("Redirecting %q => %q", r.URL.Path, redirUrl)
	http.Redirect(w, r, redirURL, 302)
}
