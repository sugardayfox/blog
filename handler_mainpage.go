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
		logger.Errorf("handleArticlesJs(): invalid sha1=%q, url='%s", sha1, url)
		serve404(w, r)
		return
	}

	jsData, expectedSha1 := getArticlesJsData()
	if sha1 != expectedSha1 {
		logger.Errorf("handleArticlesJs(): invalid value of sha1=%q, expected=%q", sha1, expectedSha1)
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

func getRecentArticles(articles []*Article, max int) []*Article {
	if max > len(articles) {
		max = len(articles)
	}
	res := make([]*Article, max, max)
	n := 0
	for i := len(articles) - 1; n < max; i-- {
		res[n] = articles[i]
		n++
	}
	return res
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
	execTemplate(w, tmpl404, model)
}

// /
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	if redirectIfNeeded(w, r) {
		return
	}

	if !isTopLevelURL(r.URL.Path) {
		serve404(w, r)
		return
	}

	articles := getCachedArticles()
	articleCount := len(articles)
	articles = getRecentArticles(articles, articleCount)

	model := struct {
		AnalyticsCode string
		Article       *Article
		Articles      []*Article
		ArticleCount  int
	}{
		AnalyticsCode: analyticsCode,
		Article:       nil, // always nil
		ArticleCount:  articleCount,
		Articles:      articles,
	}

	execTemplate(w, tmplMainPage, model)
}

// /ping
func handlePing(w http.ResponseWriter, r *http.Request) {
	servePlainText(w, "pong")
}
