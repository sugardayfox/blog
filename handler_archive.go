package main

import (
	"net/http"
)

type MonthArticle struct {
	*Article
	DisplayMonth string
}

type Year struct {
	Name     string
	Articles []MonthArticle
}

type ArticlesIndexModel struct {
	IsAdmin       bool
	AnalyticsCode string
	LogInOutURL   string
	ArticlesJsUrl string
	Article       *Article
	PostsCount    int
	Tag           string
	Years         []Year
}

func (a *MonthArticle) DisplayTitle() string {
	if a.Title != "" {
		return a.Title
	}
	return "no title"
}

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
	isAdmin := IsAdmin(r)
	articlesJsUrl := getArticlesJsUrl()
	model := ArticlesIndexModel{
		IsAdmin:       isAdmin,
		AnalyticsCode: *config.AnalyticsCode,
		LogInOutURL:   getLogInOutURL(r),
		ArticlesJsUrl: articlesJsUrl,
		PostsCount:    len(articles),
		Years:         buildYearsFromArticles(articles),
		Tag:           tag,
	}

	ExecTemplate(w, tmplArchive, model)
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
