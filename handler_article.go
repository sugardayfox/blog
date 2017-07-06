package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

// ArticleInfo describes an article
type ArticleInfo struct {
	this *Article
	next *Article
	prev *Article
	pos  int
}

func getCachedArticlesByID(articleID string) *ArticleInfo {
	articles := store.GetArticles()
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

	return getCachedArticlesByID(parts[0])
}

func makeShareHTML(r *http.Request, article *Article) string {
	title := url.QueryEscape(article.Title)
	uri := url.QueryEscape("https://" + r.Host + r.URL.String())
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

	model := struct {
		Reload         bool
		AnalyticsCode  string
		PageTitle      string
		Article        *Article
		NextArticle    *Article
		PrevArticle    *Article
		ArticlesJsURL  string
		TagsDisplay    string
		ArticleNo      int
		ArticlesCount  int
		HeaderImageURL string
		ShareHTML      template.HTML
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
		ShareHTML:     template.HTML(shareHTML),
	}

	serveTemplate(w, tmplArticle, model)
}
