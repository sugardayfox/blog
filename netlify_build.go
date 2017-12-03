package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kjk/u"
)

var (
	netlifyRedirects []*netlifyRedirect
)

type netlifyRedirect struct {
	from string
	to   string
	// valid code is 301, 302, 200, 404
	code int
}

func netlifyAddRedirect(from, to string, code int) {
	r := netlifyRedirect{
		from: from,
		to:   to,
		code: code,
	}
	netlifyRedirects = append(netlifyRedirects, &r)
}

func netlifyAddRewrite(from, to string) {
	netlifyAddRedirect(from, to, 200)
}

func netflifyAddTempRedirect(from, to string) {
	netlifyAddRedirect(from, to, 302)
}

func netflifyAddPermRedirect(from, to string) {
	netlifyAddRedirect(from, to, 301)
}

func netlifyAddStaticRedirects() {
	for from, to := range redirects {
		netflifyAddTempRedirect(from, to)
	}
}

func netlifyAddArticleRedirects() {
	for from, articleID := range articleRedirects {
		from = "/" + from
		article := store.GetArticleByID(articleID)
		u.PanicIf(article == nil, "didn't find article for id '%s'", articleID)
		to := article.URL()
		netflifyAddTempRedirect(from, to) // TODO: change to permanent
	}
	// redirect /article/:id/* => /article/:id/pretty-title
	articles := store.GetArticles(false)
	for _, article := range articles {
		from := fmt.Sprintf("/article/%s/*", article.ID)
		netflifyAddTempRedirect(from, article.URL())
	}
}

func netlifyWriteRedirects() {
	var buf bytes.Buffer
	for _, r := range netlifyRedirects {
		s := fmt.Sprintf("%s\t%s\t%d\n", r.from, r.to, r.code)
		buf.WriteString(s)
	}
	netlifyWriteFile("_redirects", buf.Bytes())
}

func mkdirForFile(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}

func copyFile(dst string, src string) error {
	err := mkdirForFile(dst)
	if err != nil {
		return err
	}
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()
	fout, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(fout, fin)
	err2 := fout.Close()
	if err != nil || err2 != nil {
		os.Remove(dst)
	}

	return err
}

func skipFile(path string) bool {
	if strings.Contains(path, ".tmpl.") {
		return true
	}
	return false
}

func dirCopyRecur(dst string, src string) (int, error) {
	fmt.Printf("dirCopyRecur: %s => %s\n", src, dst)
	nFilesCopied := 0
	dirsToVisit := []string{src}
	for len(dirsToVisit) > 0 {
		n := len(dirsToVisit)
		dir := dirsToVisit[n-1]
		dirsToVisit = dirsToVisit[:n-1]
		fileInfos, err := ioutil.ReadDir(dir)
		if err != nil {
			return nFilesCopied, err
		}
		for _, fi := range fileInfos {
			path := filepath.Join(dir, fi.Name())
			if fi.IsDir() {
				dirsToVisit = append(dirsToVisit, path)
				continue
			}
			if skipFile(path) {
				continue
			}
			dstPath := dst + path[len(src):]
			err := copyFile(dstPath, path)
			if err != nil {
				return nFilesCopied, err
			}
			nFilesCopied++
		}
	}
	return nFilesCopied, nil
}

func netlifyPath(fileName string) string {
	fileName = strings.TrimLeft(fileName, "/")
	path := filepath.Join("netlify_static", "www", fileName)
	err := mkdirForFile(path)
	u.PanicIfErr(err)
	return path
}

func netlifyWriteFile(fileName string, d []byte) {
	path := netlifyPath(fileName)
	fmt.Printf("%s\n", path)
	ioutil.WriteFile(path, d, 0644)
}

func netlifyExecTemplate(fileName string, templateName string, model interface{}) {
	path := netlifyPath(fileName)
	fmt.Printf("%s\n", path)
	var buf bytes.Buffer
	err := getTemplates().ExecuteTemplate(&buf, templateName, model)
	u.PanicIfErr(err)
	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	u.PanicIfErr(err)
}

func netlifyRequestGetFullHost() string {
	return "https://blog.kowalczyk.info"
}

func netlifyMakeShareHTML(article *Article) string {
	title := url.QueryEscape(article.Title)
	uri := netlifyRequestGetFullHost() + article.URL()
	uri = url.QueryEscape(uri)
	shareURL := fmt.Sprintf(`https://twitter.com/intent/tweet?text=%s&url=%s&via=kjk`, title, uri)
	followURL := `https://twitter.com/intent/follow?user_id=3194001`
	return fmt.Sprintf(`Hey there. You've read the whole thing. Let others know about this article by <a href="%s">sharing on Twitter</a>. <br>To be notified about new articles, <a href="%s">follow @kjk</a> on Twitter.`, shareURL, followURL)
}

func netlifyWriteArticlesArchiveForTag(tag string) {
	path := "/archives.html"
	articles := store.GetArticles(true)
	if tag != "" {
		articles = filterArticlesByTag(articles, tag, true)
		// must manually resolve conflict due to urlify
		tagInPath := tag
		if tag == "c#" {
			tagInPath = "csharp"
		} else if tag == "c++" {
			tagInPath = "cplusplus"
		}
		tagInPath = urlify(tagInPath)
		path = fmt.Sprintf("/blog/archives-by-tag-%s.html", tagInPath)
		from := "/tag/" + tag
		netlifyAddRewrite(from, path)
	}

	articlesJsURL := getArticlesJsURL()
	model := ArticlesIndexModel{
		AnalyticsCode: analyticsCode,
		ArticlesJsURL: articlesJsURL,
		PostsCount:    len(articles),
		Years:         buildYearsFromArticles(articles),
		Tag:           tag,
	}

	netlifyExecTemplate(path, tmplArchive, model)
}

func netlifyBuild() {
	// verify we're in the right directory
	_, err := os.Stat("netlify_static")
	u.PanicIfErr(err)
	outDir := filepath.Join("netlify_static", "www")
	err = os.RemoveAll(outDir)
	u.PanicIfErr(err)
	err = os.MkdirAll(outDir, 0755)
	u.PanicIfErr(err)
	nCopied, err := dirCopyRecur(outDir, "www")
	u.PanicIfErr(err)
	fmt.Printf("Copied %d files\n", nCopied)

	analyticsCode = "UA-194516-1"

	netlifyAddStaticRedirects()
	{
		// mux.HandleFunc("/contactme.html", withAnalyticsLogging(handleContactme))
		model := struct {
			RandomCookie string
		}{
			RandomCookie: randomCookie,
		}
		netlifyExecTemplate("/contactme.html", tmplContactMe, model)
	}

	{
		// 	mux.HandleFunc("/", withAnalyticsLogging(handleMainPage))
		articles := store.GetArticles(false)
		articleCount := len(articles)

		model := struct {
			AnalyticsCode string
			Article       *Article
			Articles      []*Article
			ArticleCount  int
		}{
			AnalyticsCode: analyticsCode,
			Article:       nil, // always nil
			ArticleCount:  articleCount,
			Articles:      articles,
		}

		netlifyExecTemplate("/index.html", tmplMainPage, model)
	}

	{
		// mux.HandleFunc("/atom.xml", withAnalyticsLogging(handleAtom))
		d, err := genAtomXML(true)
		u.PanicIfErr(err)
		netlifyWriteFile("/atom.xml", d)
	}

	{
		// mux.HandleFunc("/atom-all.xml", withAnalyticsLogging(handleAtomAll))
		d, err := genAtomXML(false)
		u.PanicIfErr(err)
		netlifyWriteFile("/atom-all.xml", d)
	}

	{
		// /blog/ and /kb/ are only for redirects, we only handle /article/ at this point
		articles := store.GetArticles(true)
		for _, a := range articles {
			articleInfo := getArticleInfoByID(a.ID)
			u.PanicIf(articleInfo == nil, "No article for id '%s'", a.ID)
			article := articleInfo.this
			shareHTML := netlifyMakeShareHTML(article)

			coverImage := ""
			if article.HeaderImageURL != "" {
				coverImage = netlifyRequestGetFullHost() + article.HeaderImageURL
			}

			canonicalURL := netlifyRequestGetFullHost() + article.URL()
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
				Reload:        false,
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

			path := fmt.Sprintf("/blog/%s.html", article.ID)
			netlifyExecTemplate(path, tmplArticle, model)
			netlifyAddRewrite(article.URL(), path)
		}
	}

	{
		// mux.HandleFunc("/djs/", withAnalyticsLogging(handleDjs))
		// /djs/$url
		jsData, expectedSha1 := getArticlesJsData()
		path := fmt.Sprintf("/djs/articles-%s.js", expectedSha1)
		netlifyWriteFile(path, jsData)
		from := "/djs/articles-*"
		netflifyAddTempRedirect(from, path)
	}

	{
		// mux.HandleFunc("/archives.html", withAnalyticsLogging(handleArchives))
		netlifyWriteArticlesArchiveForTag("")
		seenTags := make(map[string]bool)
		articles := store.GetArticles(false)
		for _, article := range articles {
			for _, tag := range article.Tags {
				if seenTags[tag] {
					continue
				}
				netlifyWriteArticlesArchiveForTag(tag)
				seenTags[tag] = true
			}
		}
	}

	netlifyAddArticleRedirects()
	netlifyWriteRedirects()

	/*
		mux.HandleFunc("/book/go-cookbook.html", withAnalyticsLogging(handleGoCookbook))
		mux.HandleFunc("/articles/go-cookbook.html", withAnalyticsLogging(handleGoCookbook))

		mux.HandleFunc("/sitemap.xml", withAnalyticsLogging(handleSiteMap))
		mux.HandleFunc("/software", withAnalyticsLogging(handleSoftware))
		mux.HandleFunc("/software/", withAnalyticsLogging(handleSoftware))
		mux.HandleFunc("/extremeoptimizations/", withAnalyticsLogging(handleExtremeOpt))
		mux.HandleFunc("/forum_sumatra/", withAnalyticsLogging(forumRedirect))
		mux.HandleFunc("/articles/", withAnalyticsLogging(handleArticles))
		mux.HandleFunc("/book/", withAnalyticsLogging(handleArticles))
		mux.HandleFunc("/tag/", withAnalyticsLogging(handleTag))
		mux.HandleFunc("/dailynotes-atom.xml", withAnalyticsLogging(handleNotesFeed))
		mux.HandleFunc("/dailynotes/week/", withAnalyticsLogging(handleNotesWeek))
		mux.HandleFunc("/dailynotes/tag/", withAnalyticsLogging(handleNotesTag))
		mux.HandleFunc("/dailynotes/note/", withAnalyticsLogging(handleNotesNote))
		mux.HandleFunc("/dailynotes", withAnalyticsLogging(handleNotesIndex))
		mux.HandleFunc("/worklog", handleWorkLog)
	*/
}
