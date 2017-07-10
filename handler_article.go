package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/kjk/u"
)

// ArticleInfo describes an article
type ArticleInfo struct {
	this *Article
	next *Article
	prev *Article
	pos  int
}

func getArticleInfoByID(articleID string) *ArticleInfo {
	articles := store.GetArticles(true)
	res := &ArticleInfo{}
	for i, curr := range articles {
		if curr.ID == articleID {
			if i != len(articles)-1 {
				res.next = articles[i+1]
			}
			res.this = curr
			res.pos = i
			return res
		}
		res.prev = curr
	}
	return nil
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

	return getArticleInfoByID(parts[0])
}

func makeShareHTML(r *http.Request, article *Article) string {
	title := url.QueryEscape(article.Title)
	uri := u.RequestGetFullHost(r) + r.URL.String()
	uri = url.QueryEscape(uri)
	shareURL := fmt.Sprintf(`https://twitter.com/intent/tweet?text=%s&url=%s&via=kjk`, title, uri)
	followURL := `https://twitter.com/intent/follow?user_id=3194001`
	return fmt.Sprintf(`Hey there. You've read the whole thing. Let others know about this article by <a href="%s">sharing on Twitter</a>. <br>To be notified about new articles, <a href="%s">follow @kjk</a> on Twitter.`, shareURL, followURL)
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
		logger.Noticef("handleArticle: invalid url: %s", r.URL.Path)
		serve404(w, r)
		return
	}
	article := articleInfo.this
	shareHTML := makeShareHTML(r, article)

	coverImage := ""
	if article.HeaderImageURL != "" {
		coverImage = u.RequestGetFullHost(r) + article.HeaderImageURL
	}

	canonicalURL := u.RequestGetFullHost(r) + article.URL()
	model := struct {
		Reload         bool
		AnalyticsCode  string
		PageTitle      string
		CoverImage     string
		Article        *Article
		NextArticle    *Article
		PrevArticle    *Article
		ArticlesJsURL  string
		TagsDisplay    string
		ArticleNo      int
		ArticlesCount  int
		HeaderImageURL string
		ShareHTML      template.HTML
		CanonicalURL   string
	}{
		Reload:        !flgProduction,
		AnalyticsCode: analyticsCode,
		Article:       article,
		NextArticle:   articleInfo.next,
		PrevArticle:   articleInfo.prev,
		PageTitle:     article.Title,
		CoverImage:    coverImage,
		ArticlesCount: store.ArticlesCount(),
		ArticleNo:     articleInfo.pos + 1,
		ArticlesJsURL: getArticlesJsURL(),
		ShareHTML:     template.HTML(shareHTML),
		CanonicalURL:  canonicalURL,
	}

	serveTemplate(w, tmplArticle, model)
}
