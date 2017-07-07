package main

import (
	"encoding/xml"
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
		Link:    "https://blog.kowalczyk.info/atom.xml",
		PubDate: pubTime,
	}

	for _, a := range latest {
		//id := fmt.Sprintf("tag:blog.kowalczyk.info,1999:%d", a.Id)
		e := &atom.Entry{
			Title:   a.Title,
			Link:    "https://blog.kowalczyk.info/" + a.Permalink(),
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

// SiteMapURLSet represents <urlset>
type SiteMapURLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Ns      string   `xml:"xmlns,attr"`
	URLS    []SiteMapURL
}

func makeSiteMapURLSet() *SiteMapURLSet {
	return &SiteMapURLSet{
		Ns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}
}

// SiteMapURL represents a single url
type SiteMapURL struct {
	XMLName xml.Name `xml:"url"`
	URL     string   `xml:"loc"`
	LastMod string   `xml:"lastmod"`
}

// /sitemap.xml
func handleSiteMap(w http.ResponseWriter, r *http.Request) {
	// TODO: add daily notes and static documents
	articles := getCachedArticles()
	urlset := makeSiteMapURLSet()
	var urls []SiteMapURL
	for _, article := range articles {
		articleURL := "https://" + r.Host + "/" + article.Permalink()
		uri := SiteMapURL{
			URL:     articleURL,
			LastMod: article.PublishedOn.Format("2006-01-02"),
		}
		urls = append(urls, uri)
	}
	urlset.URLS = urls

	xmlData, err := xml.MarshalIndent(urlset, " ", " ")
	if err != nil {
		serve404(w, r)
		return
	}
	d := append([]byte(xml.Header), xmlData...)
	serveXML(w, string(d))
}

// /dailynotes-atom.xml
// TODO: could cache generated xml
func handleNotesFeed(w http.ResponseWriter, r *http.Request) {
	notes := notesAllNotes
	if len(notes) > 25 {
		notes = notes[:25]
	}

	pubTime := time.Now()
	if len(notes) > 0 {
		pubTime = notes[0].Day
	}

	feed := &atom.Feed{
		Title:   "Krzysztof Kowalczyk daily notes",
		Link:    "https://blog.kowalczyk.info/dailynotes-atom.xml",
		PubDate: pubTime,
	}

	for _, n := range notes {
		//id := fmt.Sprintf("tag:blog.kowalczyk.info,1999:%d", a.Id)
		title := n.Title
		if title == "" {
			title = n.ID
		}
		html := `<pre>` + n.HTMLBody + `</pre>`
		e := &atom.Entry{
			Title:   title,
			Link:    "https://blog.kowalczyk.info/" + n.URL,
			Content: html,
			PubDate: n.Day,
		}
		feed.AddEntry(e)
	}

	s, err := feed.GenXml()
	if err != nil {
		s = []byte("Failed to generate XML feed")
	}

	w.Write(s)
}
