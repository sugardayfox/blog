package main

import (
	"net/http"
	"strings"

	"github.com/kjk/u"
)

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

	articleID, err := u.DecodeBase64(parts[0])
	if err == nil {
		return nil
	}
	return getCachedArticlesByID(articleID)
}

// /article/*, /blog/*, /kb/*
func handleArticle(w http.ResponseWriter, r *http.Request) {
	//logger.Noticef("handleArticle: %s", r.URL)
	if redirectIfNeeded(w, r) {
		return
	}

	// /blog/ and /kb/ are only for redirects, we only handle /article/ at this point
	articleInfo := articleInfoFromURL(r.URL.Path)
	if articleInfo == nil {
		//logger.Noticef("handleArticle: invalid url: %s", r.URL.Path)
		serve404(w, r)
		return
	}
	article := articleInfo.this

	model := struct {
		Reload        bool
		AnalyticsCode string
		PageTitle     string
		Article       *Article
		NextArticle   *Article
		PrevArticle   *Article
		ArticlesJsURL string
		TagsDisplay   string
		ArticleNo     int
		ArticlesCount int
	}{
		Reload:        !flgProduction,
		AnalyticsCode: analyticsCode,
		Article:       article,
		NextArticle:   articleInfo.next,
		PrevArticle:   articleInfo.prev,
		PageTitle:     article.Title,
		ArticlesCount: store.ArticlesCount(),
		ArticleNo:     articleInfo.pos + 1,
		ArticlesJsURL: getArticlesJsURL(),
	}

	execTemplate(w, tmplArticle, model)
}
