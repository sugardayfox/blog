package main

import (
	"net/http"
	"strconv"
	"strings"
)

// url is in the form: ${sha1}.js
func handleArticlesJs(w http.ResponseWriter, r *http.Request, url string) {
	sha1 := url[:len(url)-len(".js")]
	if len(sha1) != 40 {
		serve404(w, r)
		return
	}

	jsData, expectedSha1 := getArticlesJsData()
	if sha1 != expectedSha1 {
		// this might happen due to caching and stale url, return the old value
	}

	w.Header().Set("Content-Type", "text/javascript")
	// cache non-admin version by setting max age 1 year into the future
	// http://betterexplained.com/articles/how-to-optimize-your-site-with-http-caching/
	w.Header().Set("Cache-Control", "max-age=31536000, public")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsData)))
	w.Write(jsData)
}

// /djs/$url
func handleDjs(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/djs/"):]
	if strings.HasPrefix(url, "articles-") {
		handleArticlesJs(w, r, url[len("articles-"):])
		return
	}
	serve404(w, r)
}

func serve404(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	model := struct {
		URL string
	}{
		URL: uri,
	}
	serveTemplate(w, tmpl404, model)
}
