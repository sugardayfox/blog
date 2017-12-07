package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	_ "net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kjk/u"
)

var (
	store         *ArticlesStore
	analyticsCode = "UA-194516-1"

	flgNewArticleTitle string
)

func parseCmdLineFlags() {
	flag.StringVar(&flgNewArticleTitle, "newarticle", "", "create a new article")
	flag.Parse()
}

func findUniqueArticleID(articles []*Article) string {
	existingIDs := make(map[string]bool)
	for _, a := range articles {
		existingIDs[a.ID] = true
	}

	for i := 1; i < 10000; i++ {
		s := u.EncodeBase64(i)
		if !existingIDs[s] {
			return strconv.Itoa(i)
		}
	}
	u.PanicIf(true, "couldn't find unique article id")
	return ""
}

func genNewArticle(title string) {
	fmt.Printf("genNewArticle: %q\n", title)
	store, err := NewArticlesStore()
	if err != nil {
		log.Fatalf("NewStore() failed with %s", err)
	}
	newID := findUniqueArticleID(store.articlesWithDrafts)
	t := time.Now()
	dir := "articles"
	yyyy := fmt.Sprintf("%04d", t.Year())
	month := t.Month()
	sanitizedTitle := sanitizeForFile(title)
	name := fmt.Sprintf("%02d-%s.md", month, sanitizedTitle)
	fmt.Printf("new id: %s, name: %s\n", newID, name)
	path := filepath.Join(dir, yyyy, name)
	s := fmt.Sprintf(`---
Id: %s
Title: %s
Date: %s
Format: Markdown
---`, newID, title, t.Format(time.RFC3339))
	for i := 1; i < 10; i++ {
		if !u.PathExists(path) {
			break
		}
		name = fmt.Sprintf("%02d-%s-%d.md", month, sanitizedTitle, i)
		path = filepath.Join(dir, yyyy, name)
	}
	u.PanicIf(u.PathExists(path))
	fmt.Printf("path: %s\n", path)
	u.CreateDirForFileMust(path)
	ioutil.WriteFile(path, []byte(s), 0644)
}

func loadArticles() {
	var err error
	if store, err = NewArticlesStore(); err != nil {
		log.Fatalf("NewStore() failed with %s", err)
	}
	articles := store.GetArticles(true)
	articlesJs, articlesJsSha1 = buildArticlesJSON(articles)
}

// https://caddyserver.com/tutorial/caddyfile
var caddyProlog = `
localhost:8080
root netlify_static
`

func writeCaddyConfig() {
	path := filepath.Join("Caddyfile")
	f, err := os.Create(path)
	u.PanicIfErr(err)
	defer f.Close()

	_, err = f.Write([]byte(caddyProlog))
	u.PanicIfErr(err)
	var s string
	for _, r := range netlifyRedirects {
		if r.code == 200 {
			s = fmt.Sprintf("rewrite %s %s %d\n", r.from, r.to, r.code)
		} else {
			s = fmt.Sprintf("redir %s %s %d\n", r.from, r.to, r.code)
		}
		_, err = io.WriteString(f, s)
		u.PanicIfErr(err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	parseCmdLineFlags()

	if flgNewArticleTitle != "" {
		genNewArticle(flgNewArticleTitle)
		return
	}

	notesGenIDIfNecessary()
	loadTemplates()
	loadArticles()
	readRedirects()
	netlifyBuild()
	writeCaddyConfig()
}
