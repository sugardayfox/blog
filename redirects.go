package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/kjk/u"
)

var redirects = map[string]string{
	"/index.html":                                   "/",
	"/blog":                                         "/",
	"/blog/":                                        "/",
	"/kb/serialization-in-c#.html":                  "/article/Serialization-in-C.html",
	"/extremeoptimizations":                         "/extremeoptimizations/index.html",
	"/extremeoptimizations/":                        "/extremeoptimizations/index.html",
	"/feed/rss2/atom.xml":                           "/atom.xml",
	"/feed/rss2/":                                   "/atom.xml",
	"/feed/rss2":                                    "/atom.xml",
	"/feed/":                                        "/atom.xml",
	"/feed":                                         "/atom.xml",
	"/feedburner.xml":                               "/atom.xml",
	"/articles/cocoa-objectivec-reference.html":     "/articles/cocoa-reference.html",
	"/forum_sumatra":                                "https://forum.sumatrapdfreader.org/",
	"/google6dba371684d43cd6.html":                  "/static/google6dba371684d43cd6.html",
	"/software/15minutes/index.html":                "/software/15minutes.html",
	"/software/15minutes/":                          "/software/15minutes.html",
	"/software/fofou":                               "/software/fofou/index.html",
	"/software/patheditor":                          "/software/patheditor/for-windows.html",
	"/software/patheditor/":                         "/software/patheditor/for-windows.html",
	"/software/scdiff/":                             "/software/scdiff.html",
	"/software/scdiff/index.html":                   "/software/scdiff.html",
	"/free-pdf-reader.html":                         "https://www.sumatrapdfreader.org/free-pdf-reader.html",
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

var articleRedirects = make(map[string]string)
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
		u.PanicIf(len(parts) != 2, "malformed article_redirects.txt, len(parts) = %d (!2)", len(parts))
		idStr := parts[0]
		url := strings.TrimSpace(parts[1])
		idNum, err := strconv.Atoi(idStr)
		u.PanicIfErr(err, "malformed line in article_redirects.txt. Line:\n%s\n", l)
		id := u.EncodeBase64(idNum)
		a := store.GetArticleByID(id)
		if a != nil {
			articleRedirects[url] = id
			continue
		}
		//fmt.Printf("skipping redirect '%s' because article with id %d no longer present\n", string(l), id)
	}
}
