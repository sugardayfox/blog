package main

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

const (
	formatHTML     = 0
	formatMarkdown = 1
	formatText     = 2

	formatUnknown = -1
)

// Article describes a single article
type Article struct {
	ID          int
	PublishedOn time.Time
	Title       string
	Tags        []string
	Format      int
	Hidden      bool // is it hidden from main timeline?
	Path        string
	Body        []byte
	BodyHTML    string

	HTMLBody     template.HTML
	DisplayMonth string
}

// same format as Format* constants
var formatNames = []string{"Html", "Markdown", "Text"}

func validFormat(format int) bool {
	return format >= formatHTML && format <= formatText
}

func urlForTag(tag string) string {
	// TODO: url-quote the first tag
	return fmt.Sprintf(`<a href="/tag/%s" class="taglink">%s</a>`, tag, tag)
}

// FormatNameToID return id of a format
func FormatNameToID(name string) int {
	for i, formatName := range formatNames {
		if strings.EqualFold(name, formatName) {
			return i
		}
	}
	return formatUnknown
}

// Permalink returns article's permalink
func (a *Article) Permalink() string {
	return "article/" + shortenID(a.ID) + "/" + urlify(a.Title) + ".html"
}

// TagsDisplay returns tags as html
func (a *Article) TagsDisplay() template.HTML {
	arr := make([]string, 0)
	for _, tag := range a.Tags {
		arr = append(arr, urlForTag(tag))
	}
	s := strings.Join(arr, ", ")
	return template.HTML(s)
}

// PublishedOnShort is a short version of date
func (a *Article) PublishedOnShort() string {
	return a.PublishedOn.Format("Jan 2 2006")
}
