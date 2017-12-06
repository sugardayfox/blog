package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chilts/sid"
	"github.com/kjk/betterguid"
	"github.com/kjk/u"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"
	atom "github.com/thomas11/atomgenerator"
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

func copyAndSortArticles(articles []*Article) []*Article {
	n := len(articles)
	res := make([]*Article, n, n)
	copy(res, articles)
	sort.Slice(res, func(i, j int) bool {
		return res[j].PublishedOn.After(res[i].PublishedOn)
	})
	return res
}

func genAtomXML(excludeNotes bool) ([]byte, error) {
	articles := store.GetArticles(false)
	if excludeNotes {
		articles = filterArticlesByTag(articles, "note", false)
	}
	articles = copyAndSortArticles(articles)
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
			Link:    "https://blog.kowalczyk.info" + a.URL(),
			Content: a.BodyHTML,
			PubDate: a.PublishedOn,
		}
		feed.AddEntry(e)
	}

	return feed.GenXml()
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
	netlifyAddRewrite("/contactme.html", "/static/contactme-netlify.html")
	netlifyAddRewrite("/favicon.ico", "/static/favicon.ico")
	netlifyAddRewrite("/articles/", "/static/documents.html")
	netlifyAddRewrite("/articles/index.html", "/static/documents.html")
	netlifyAddRewrite("/book/", "/static/documents.html")
	netflifyAddTempRedirect("/book/*", "/articles/:splat")

	netlifyExecTemplate("/static/documents.html", tmplDocuments, nil)

	{
		// url: /book/go-cookbook.html
		model := struct {
			InProduction bool
		}{
			InProduction: true,
		}
		netlifyExecTemplate("/book/go-cookbook.html", tmplGoCookBook, model)
		netlifyAddRewrite("/articles/go-cookbook.html", "/book/go-cookbook.html")
	}

	{
		// 	mux.HandleFunc("/", handleMainPage)
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
		// /atom.xml
		d, err := genAtomXML(true)
		u.PanicIfErr(err)
		netlifyWriteFile("/atom.xml", d)
	}

	{
		// /atom-all.xml
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
		// mux.HandleFunc("/djs/", handleDjs)
		// /djs/$url
		jsData, expectedSha1 := getArticlesJsData()
		path := fmt.Sprintf("/djs/articles-%s.js", expectedSha1)
		netlifyWriteFile(path, jsData)
		from := "/djs/articles-*"
		netflifyAddTempRedirect(from, path)
	}

	{
		// mux.HandleFunc("/archives.html", handleArchives)
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

	{
		// mux.HandleFunc("/dailynotes", handleNotesIndex)
		weekStart := notesWeekStarts[0]
		notes := notesWeekStartDayToNotes[weekStart]
		var nextWeek string
		if len(notesWeekStarts) > 1 {
			nextWeek = notesWeekStarts[1]
		}
		model := &modelNotesForWeek{
			Notes:         notes,
			TagCounts:     notesTagCounts,
			TotalNotes:    nTotalNotes,
			WeekStartDay:  weekStart,
			AnalyticsCode: analyticsCode,
			NextWeek:      nextWeek,
		}
		netlifyExecTemplate("/dailynotes.html", tmplNotesWeek, model)
		netlifyAddRewrite("/dailynotes", "dailynotes.html")
	}

	{
		// /dailynotes/week/${day} : week starting with a given day
		for weekStart, notes := range notesWeekStartDayToNotes {
			u.PanicIf(len(notes) == 0, "no notes for week '%s'", weekStart)
			var nextWeek, prevWeek string
			for idx, ws := range notesWeekStarts {
				if ws != weekStart {
					continue
				}
				if idx > 0 {
					prevWeek = notesWeekStarts[idx-1]
				}
				lastIdx := len(notesWeekStarts) - 1
				if idx+1 <= lastIdx {
					nextWeek = notesWeekStarts[idx+1]
				}
				break
			}
			model := &modelNotesForWeek{
				Notes:         notes,
				TagCounts:     notesTagCounts,
				WeekStartDay:  weekStart,
				NextWeek:      nextWeek,
				PrevWeek:      prevWeek,
				AnalyticsCode: analyticsCode,
			}
			path := fmt.Sprintf("/dailynotes-week-%s.html", weekStart)
			netlifyExecTemplate(path, tmplNotesWeek, model)
			from := "/dailynotes/week/" + weekStart
			netlifyAddRewrite(from, path)
		}
	}

	{
		// /dailynotes/note/${id}-${title}
		for _, aNote := range notesAllNotes {
			weekStartTime := calcWeekStart(aNote.Day)
			weekStartDay := weekStartTime.Format("2006-01-02")
			model := struct {
				WeekStartDay  string
				Note          *note
				AnalyticsCode string
			}{
				WeekStartDay:  weekStartDay,
				Note:          aNote,
				AnalyticsCode: analyticsCode,
			}
			path := fmt.Sprintf("/dailynotes-note-%s.html", aNote.ID)
			netlifyExecTemplate(path, tmplNotesNote, model)
			from := fmt.Sprintf("/dailynotes/note/%s", aNote.ID)
			if aNote.Title != "" {
				from += "-" + urlify(aNote.Title)
			}
			netlifyAddRewrite(from, path)
		}
	}

	{
		// /dailynotes/tag/${tag}
		seenTags := make(map[string]bool)
		for tag, notes := range notesTagToNotes {
			u.PanicIf(len(notes) == 0, "no notes for tag '%s'", tag)
			model := struct {
				Notes         []*note
				TagCounts     []tagWithCount
				Tag           string
				AnalyticsCode string
			}{
				Notes:         notes,
				TagCounts:     notesTagCounts,
				Tag:           tag,
				AnalyticsCode: analyticsCode,
			}
			// TODO: this tag can be
			tag2 := urlify(tag)
			u.PanicIf(seenTags[tag2], "already seen tag: '%s' '%s'", tag, tag2)
			path := fmt.Sprintf("/dailynotes-tag-%s.html", tag2)
			netlifyExecTemplate(path, tmplNotesTag, model)
			from := fmt.Sprintf("/dailynotes/tag/%s", tag)
			netlifyAddRewrite(from, path)
		}
	}

	{
		// /dailynotes-atom.xml
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
			e := &atom.Entry{
				Title:   title,
				Link:    "https://blog.kowalczyk.info/" + n.URL,
				Content: string(n.HTMLBody),
				PubDate: n.Day,
			}
			feed.AddEntry(e)
		}

		data, err := feed.GenXml()
		u.PanicIfErr(err)
		netlifyWriteFile("/dailynotes-atom.xml", data)
	}

	{
		// /sitemap.xml
		data, err := genSiteMap("https://blog.kowalczyk.info")
		u.PanicIfErr(err)
		netlifyWriteFile("/sitemap.xml", data)
	}

	{
		// /tools/generate-unique-id
		idXid := xid.New()
		idKsuid := ksuid.New()

		t := time.Now().UTC()
		entropy := rand.New(rand.NewSource(t.UnixNano()))
		idUlid := ulid.MustNew(ulid.Timestamp(t), entropy)
		betterGUID := betterguid.New()
		uuid := uuid.NewV4()

		flake := sonyflake.NewSonyflake(sonyflake.Settings{})
		sfid, err := flake.NextID()
		sfidstr := fmt.Sprintf("%x", sfid)
		if err != nil {
			sfidstr = err.Error()
		}

		model := struct {
			Xid           string
			Ksuid         string
			Ulid          string
			BetterGUID    string
			Sonyflake     string
			Sid           string
			UUIDv4        string
			AnalyticsCode string
		}{
			Xid:           idXid.String(),
			Ksuid:         idKsuid.String(),
			Ulid:          idUlid.String(),
			BetterGUID:    betterGUID,
			Sonyflake:     sfidstr,
			Sid:           sid.Id(),
			UUIDv4:        uuid.String(),
			AnalyticsCode: analyticsCode,
		}

		// make sure /tools/generate-unique-id is served as html
		path := "/tools/generate-unique-id.html"
		netlifyExecTemplate(path, tmplGenerateUniqueID, model)
		netlifyAddRewrite("/tools/generate-unique-id", path)
	}

	// no longer care about /worklog

	netlifyAddArticleRedirects()
	netlifyWriteRedirects()

	/*
		mux.HandleFunc("/extremeoptimizations/", handleExtremeOpt)
	*/
}
