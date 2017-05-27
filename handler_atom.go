package main

import (
	"net/http"
	"time"

	atom "github.com/thomas11/atomgenerator"
)

func handleAtomHelp(w http.ResponseWriter, r *http.Request, excludeNotes bool) {
	articles := getCachedArticles()
	if excludeNotes {
		articles = filterArticlesByTag(articles, "note", false)
	}
	n := 25
	if n > len(articles) {
		n = len(articles)
	}

	latest := make([]*Article, n, n)
	size := len(articles)
	for i := 0; i < n; i++ {
		latest[i] = articles[size-1-i]
	}

	pubTime := time.Now()
	if len(articles) > 0 {
		pubTime = articles[0].PublishedOn
	}

	feed := &atom.Feed{
		Title:   "Krzysztof Kowalczyk blog",
		Link:    "http://blog.kowalczyk.info/atom.xml",
		PubDate: pubTime,
	}

	for _, a := range latest {
		//id := fmt.Sprintf("tag:blog.kowalczyk.info,1999:%d", a.Id)
		e := &atom.Entry{
			Title:   a.Title,
			Link:    "http://blog.kowalczyk.info/" + a.Permalink(),
			Content: a.BodyHTML,
			PubDate: a.PublishedOn,
		}
		feed.AddEntry(e)
	}

	s, err := feed.GenXml()
	if err != nil {
		s = []byte("Failed to generate XML feed")
	}

	w.Write(s)
}

// /atom-all.xml
func handleAtomAll(w http.ResponseWriter, r *http.Request) {
	handleAtomHelp(w, r, false)
}

// /atom.xml
func handleAtom(w http.ResponseWriter, r *http.Request) {
	handleAtomHelp(w, r, true)
}
