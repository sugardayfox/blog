package main

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/kjk/u"
)

var articlesCache ArticlesCache

type ArticlesCache struct {
	sync.Mutex
	articles       []*Article
	articlesJs     []byte
	articlesJsSha1 string
}

func appendJsonMarshalled(buf *bytes.Buffer, val interface{}) {
	if data, err := json.Marshal(val); err != nil {
		logger.Errorf("json.Marshal() of %v failed with %s", val, err)
	} else {
		buf.Write(data)
	}
}

// TODO: I only use it for tag cloud, could just send info about tags directly
func buildArticlesJson(articles []*Article) ([]byte, string) {
	var buf bytes.Buffer
	buf.WriteString("var __articles_json = ")
	n := len(articles)
	vals := make([]interface{}, n, n)
	n = 0
	for i := len(articles) - 1; i >= 0; i-- {
		a := articles[i]
		val := make([]interface{}, 1, 1)
		val[0] = a.Tags
		vals[n] = val
		n += 1
	}
	appendJsonMarshalled(&buf, vals)
	buf.WriteString("; articlesJsonLoaded(__articles_json);")
	jsData := buf.Bytes()
	sha1 := u.Sha1HexOfBytes(jsData)
	//logger.Noticef("buildArticlesJson(): len(jsData)=%d, sha1=%s", len(jsData), sha1)
	return jsData, sha1
}

func buildArticlesCache() {
	u.PanicIf(articlesCache.articles != nil)
	articles := store.GetArticles()
	articlesCache.articles = articles
	articlesCache.articlesJs, articlesCache.articlesJsSha1 = buildArticlesJson(articles)
}

func getArticlesJsUrl() string {
	sha1 := articlesCache.articlesJsSha1
	return "/djs/articles-" + sha1 + ".js"
}

func getArticlesJsData() ([]byte, string) {
	return articlesCache.articlesJs, articlesCache.articlesJsSha1
}

func getCachedArticles() []*Article {
	return articlesCache.articles
}

type ArticleInfo struct {
	this *Article
	next *Article
	prev *Article
	pos  int
}

func getCachedArticlesById(articleId int) *ArticleInfo {
	articles := store.GetArticles()
	res := &ArticleInfo{}
	for i, curr := range articles {
		if curr.ID == articleId {
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
