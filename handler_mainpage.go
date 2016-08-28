package main

import (
	"net/http"
	"strconv"
	"strings"
)

// url is in the form: $sha1.js
func handleArticlesJs(w http.ResponseWriter, r *http.Request, url string) {
	sha1 := url[:len(url)-len(".js")]
	if len(sha1) != 40 {
		logger.Errorf("handleArticlesJs(): invalid sha1=%q, url='%s", sha1, url)
		panic("invalid sha1")
	}

	jsData, expectedSha1 := getArticlesJsData()
	if sha1 != expectedSha1 {
		logger.Errorf("handleArticlesJs(): invalid value of sha1=%q, expected=%q", sha1, expectedSha1)
		panic("invalid value of sha1")
	}

	w.Header().Set("Content-Type", "text/javascript")
	// cache non-admin version by setting max age 1 year into the future
	// http://betterexplained.com/articles/how-to-optimize-your-site-with-http-caching/
	if !IsAdmin(r) {
		w.Header().Set("Cache-Control", "max-age=31536000, public")
	}
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
	http.NotFound(w, r)
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

// /
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	if redirectIfNeeded(w, r) {
		return
	}

	if !isTopLevelURL(r.URL.Path) {
		http.NotFound(w, r)
		return
	}

	isAdmin := IsAdmin(r)
	articles := getCachedArticles()
	articleCount := len(articles)
	articles = getRecentArticles(articles, articleCount)

	model := struct {
		IsAdmin       bool
		AnalyticsCode string
		Article       *Article
		Articles      []*Article
		ArticleCount  int
		LogInOutURL   string
	}{
		IsAdmin:       isAdmin,
		AnalyticsCode: *config.AnalyticsCode,
		Article:       nil, // always nil
		ArticleCount:  articleCount,
		Articles:      articles,
		LogInOutURL:   getLogInOutURL(r),
	}

	ExecTemplate(w, tmplMainPage, model)
}
