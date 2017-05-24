package main

import (
	"bytes"
	"encoding/json"
	"sync"
)

var articlesCache ArticlesCache

// ArticlesCache describes a cache of articles
type ArticlesCache struct {
	sync.Mutex
	articles       []*Article
	articlesJs     []byte
	articlesJsSha1 string
}

func appendJSONMarshalled(buf *bytes.Buffer, val interface{}) {
	if data, err := json.Marshal(val); err != nil {
		logger.Errorf("json.Marshal() of %v failed with %s", val, err)
	} else {
		buf.Write(data)
	}
}

// TODO: I only use it for tag cloud, could just send info about tags directly
func buildArticlesJSON(articles []*Article) ([]byte, string) {
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
		n++
	}
	appendJSONMarshalled(&buf, vals)
	buf.WriteString("; articlesJsonLoaded(__articles_json);")
	jsData := buf.Bytes()
	sha1 := sha1HexOfBytes(jsData)
	//logger.Noticef("buildArticlesJson(): len(jsData)=%d, sha1=%s", len(jsData), sha1)
	return jsData, sha1
}

func getArticlesJsURL() string {
	sha1 := articlesCache.articlesJsSha1
	return "/djs/articles-" + sha1 + ".js"
}

func getArticlesJsData() ([]byte, string) {
	return articlesCache.articlesJs, articlesCache.articlesJsSha1
}

func getCachedArticles() []*Article {
	return articlesCache.articles
}

// ArticleInfo describes an article
type ArticleInfo struct {
	this *Article
	next *Article
	prev *Article
	pos  int
}

func getCachedArticlesByID(articleID int) *ArticleInfo {
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
