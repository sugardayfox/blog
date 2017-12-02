package main

import (
	"encoding/xml"
	"net/http"
	"path"
	"time"

	"github.com/kjk/blog/pkg/notes"
)

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
	XMLName      xml.Name `xml:"url"`
	URL          string   `xml:"loc"`
	LastModified string   `xml:"lastmod"`
}

// There are more static pages, but those are the important ones
var staticURLS = []string{
	"/book/go-cookbook.html",
	"/articles/cbz-cbr-comic-book-reader-viewer-for-windows.html",
	"/articles/chm-reader-viewer-for-windows.html",
	"/articles/mobi-ebook-reader-viewer-for-windows.html",
	"/articles/epub-ebook-reader-viewer-for-windows.html",
	"/articles/where-to-get-free-ebooks-epub-mobi.html",
	"/software/",
	"static/documents.html",
	"/dailynotes",
}

// /sitemap.xml
func handleSiteMap(w http.ResponseWriter, r *http.Request) {
	articles := store.GetArticles(true)
	urlset := makeSiteMapURLSet()
	var urls []SiteMapURL
	for _, article := range articles {
		pageURL := "https://" + path.Join(r.Host, article.URL())
		uri := SiteMapURL{
			URL:          pageURL,
			LastModified: article.UpdatedOn.Format("2006-01-02"),
		}
		urls = append(urls, uri)
	}

	now := time.Now()
	for _, staticURL := range staticURLS {
		pageURL := "https://" + path.Join(r.Host, staticURL)
		uri := SiteMapURL{
			URL:          pageURL,
			LastModified: now.Format("2006-01-02"),
		}
		urls = append(urls, uri)
	}

	for _, note := range notes.NotesAllNotes {
		pageURL := "https://" + path.Join(r.Host, note.URL)
		uri := SiteMapURL{
			URL:          pageURL,
			LastModified: note.Day.Format("2006-01-02"),
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
