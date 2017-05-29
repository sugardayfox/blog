package main

import (
	"net/http"
)

// MonthArticle combines article and a month
type MonthArticle struct {
	*Article
	DisplayMonth string
}

// Year describes articles in a given year
type Year struct {
	Name     string
	Articles []MonthArticle
}

// ArticlesIndexModel describes index of articles
type ArticlesIndexModel struct {
	AnalyticsCode string
	ArticlesJsURL string
	Article       *Article
	PostsCount    int
	Tag           string
	Years         []Year
}

// DisplayTitle returns a title for an article
func (a *MonthArticle) DisplayTitle() string {
	if a.Title != "" {
		return a.Title
	}
	return "no title"
}

// NewYear creates a new Year
func NewYear(name string) *Year {
	return &Year{Name: name, Articles: make([]MonthArticle, 0)}
}

func buildYearsFromArticles(articles []*Article) []Year {
	res := make([]Year, 0)
	var currYear *Year
	var currMonthName string
	n := len(articles)
	for i := n - 1; i >= 0; i-- {
		a := articles[i]
		yearName := a.PublishedOn.Format("2006")
		if currYear == nil || currYear.Name != yearName {
			if currYear != nil {
				res = append(res, *currYear)
			}
			currYear = NewYear(yearName)
			currMonthName = ""
		}
		ma := MonthArticle{Article: a}
		monthName := a.PublishedOn.Format("01")
		if monthName != currMonthName {
			ma.DisplayMonth = a.PublishedOn.Format("January 2")
		} else {
			ma.DisplayMonth = a.PublishedOn.Format("2")
		}
		currMonthName = monthName
		currYear.Articles = append(currYear.Articles, ma)
	}
	if currYear != nil {
		res = append(res, *currYear)
	}
	return res
}

func filterArticlesByTag(articles []*Article, tag string, include bool) []*Article {
	res := make([]*Article, 0)
	for _, a := range articles {
		hasTag := false
		for _, t := range a.Tags {
			if tag == t {
				hasTag = true
				break
			}
		}
		if include && hasTag {
			res = append(res, a)
		} else if !include && !hasTag {
			res = append(res, a)
		}
	}
	return res
}

func showArchiveArticles(w http.ResponseWriter, r *http.Request, articles []*Article, tag string) {
	articlesJsURL := getArticlesJsURL()
	model := ArticlesIndexModel{
		AnalyticsCode: analyticsCode,
		ArticlesJsURL: articlesJsURL,
		PostsCount:    len(articles),
		Years:         buildYearsFromArticles(articles),
		Tag:           tag,
	}

	execTemplate(w, tmplArchive, model)
}

func showArchivePage(w http.ResponseWriter, r *http.Request, tag string) {
	articles := getCachedArticles()
	if tag != "" {
		articles = filterArticlesByTag(articles, tag, true)
	}
	showArchiveArticles(w, r, articles, tag)
}

// /tag/${tag}
func handleTag(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Path[len("/tag/"):]
	showArchivePage(w, r, tag)
}

// /archives.html
func handleArchives(w http.ResponseWriter, r *http.Request) {
	showArchivePage(w, r, "")
}
