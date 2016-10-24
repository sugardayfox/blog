package main

import (
	"html/template"
	"net/http"
	"strings"
)

// DisplayArticle represents an article to display
type DisplayArticle struct {
	*Article
	HTMLBody template.HTML
}

// PublishedOnShort is a short version of date
func (a *DisplayArticle) PublishedOnShort() string {
	return a.PublishedOn.Format("Jan 2 2006")
}

func articleInfoFromURL(uri string) *ArticleInfo {
	if strings.HasPrefix(uri, "/") {
		uri = uri[1:]
	}
	if !strings.HasPrefix(uri, "article/") {
		return nil
	}
	// we expect /article/$shortId/$url
	parts := strings.SplitN(uri[len("article/"):], "/", 2)
	if len(parts) != 2 {
		return nil
	}

	articleID := UnshortenID(parts[0])
	return getCachedArticlesByID(articleID)
}

// /article/*, /blog/*, /kb/*
func handleArticle(w http.ResponseWriter, r *http.Request) {
	//logger.Noticef("handleArticle: %s", r.URL)
	if redirectIfNeeded(w, r) {
		return
	}
	isAdmin := IsAdmin(r)

	// /blog/ and /kb/ are only for redirects, we only handle /article/ at this point
	uri := r.URL.Path
	articleInfo := articleInfoFromURL(r.URL.Path)
	if articleInfo == nil {
		logger.Noticef("handleArticle: invalid url: %s\n", uri)
		http.NotFound(w, r)
		return
	}
	article := articleInfo.this
	displayArticle := &DisplayArticle{Article: article}
	msgHTML := article.GetHTMLStr()
	displayArticle.HTMLBody = template.HTML(msgHTML)

	model := struct {
		IsAdmin       bool
		Reload        bool
		AnalyticsCode string
		PageTitle     string
		Article       *DisplayArticle
		NextArticle   *Article
		PrevArticle   *Article
		LogInOutURL   string
		ArticlesJsURL string
		TagsDisplay   string
		ArticleNo     int
		ArticlesCount int
	}{
		IsAdmin:       isAdmin,
		Reload:        !inProduction,
		AnalyticsCode: *config.AnalyticsCode,
		LogInOutURL:   getLogInOutURL(r),
		Article:       displayArticle,
		NextArticle:   articleInfo.next,
		PrevArticle:   articleInfo.prev,
		PageTitle:     article.Title,
		ArticlesCount: store.ArticlesCount(),
		ArticleNo:     articleInfo.pos + 1,
		ArticlesJsURL: getArticlesJsURL(),
	}

	ExecTemplate(w, tmplArticle, model)
}
