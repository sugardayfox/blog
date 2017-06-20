package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kjk/u"
	"github.com/kr/fs"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

const (
	formatHTML     = 0
	formatMarkdown = 1
	formatText     = 2

	formatUnknown = -1
)

// Article describes a single article
type Article struct {
	ID             string
	PublishedOn    time.Time
	Title          string
	Tags           []string
	Format         int
	Hidden         bool // is it hidden from main timeline?
	Path           string
	Body           []byte
	BodyHTML       string
	HeaderImageURL string

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
	return "article/" + a.ID + "/" + urlify(a.Title) + ".html"
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

// ArticlesStore is a store for articles
type ArticlesStore struct {
	articles    []*Article
	idToArticle map[string]*Article
	dirsToWatch []string
}

func isSepLine(s string) bool {
	return strings.HasPrefix(s, "-----")
}

func parseTags(s string) []string {
	tags := strings.Split(s, ",")
	var res []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		tag = strings.ToLower(tag)
		// skip the tag I use in quicknotes.io to tag notes for the blog
		if tag == "for-blog" || tag == "published" || tag == "draft" {
			continue
		}
		res = append(res, tag)
	}
	return res
}

func parseFormat(s string) int {
	s = strings.ToLower(s)
	switch s {
	case "html":
		return formatHTML
	case "markdown", "md":
		return formatMarkdown
	case "text":
		return formatText
	default:
		return formatUnknown
	}
}

func parseDate(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(s, "2006-01-02")
	if err == nil {
		return t, nil
	}
	// TODO: more formats?
	return time.Now(), err
}

func extractMetadataValue(d []byte, prefix string) ([]byte, string) {
	if !bytes.HasPrefix(d, []byte(prefix)) {
		return d, ""
	}
	d = bytes.TrimPrefix(d, []byte(prefix))
	eolIdx := bytes.IndexByte(d, '\n')
	if eolIdx == -1 {
		return []byte{}, string(d)
	}
	val := d[:eolIdx]
	d = d[eolIdx+1:]
	return d, strings.TrimSpace(string(val))
}

// a note can have additional metadata at the beginning in the form of:
// @{name} ${value}\n
// We extract this metadata and put the relevant info in article
func extractAdditionalMetadata(d []byte, article *Article) ([]byte, error) {
	var fileName string
	oneMore := true
	for oneMore {
		oneMore = false
		d, fileName = extractMetadataValue(d, "@header-image")
		if fileName != "" {
			oneMore = true
			if fileName[0] != '/' {
				fileName = "/" + fileName
			}
			path := filepath.Join("www", fileName)
			if !u.FileExists(path) {
				return d, fmt.Errorf("File '%s' for @header-image doesn't exist", path)
			}
			//fmt.Printf("Found HeaderImageURL: %s\n", fileName)
			article.HeaderImageURL = fileName
			continue
		}
	}
	return d, nil
}

// might return nil if article is meant to be skipped (deleted or draft)
func readArticle(path string) (*Article, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	a := &Article{}
	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		l = strings.TrimSpace(l)
		if isSepLine(l) {
			break
		}
		parts := strings.SplitN(l, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Unexpected line: %q", l)
		}
		k := strings.ToLower(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "deleted":
			return nil, nil
		case "draft":
			if flgProduction {
				return nil, nil
			}
		case "id":
			// we handle 2 types of ids
			// blog posts from articles/ directory have integer id
			// blog posts imported from quicknotes (articles/from-quicknotes/)
			// have id that are strings
			id, err := strconv.Atoi(v)
			if err == nil {
				a.ID = u.EncodeBase64(id)
			} else {
				a.ID = strings.TrimSpace(v)
			}
			a.Path = path
		case "title":
			a.Title = v
		case "tags":
			a.Tags = parseTags(v)
		case "format":
			f := parseFormat(v)
			if f == formatUnknown {
				return nil, fmt.Errorf("%q is not a valid format", v)
			}
			a.Format = f
		case "date", "createdat":
			a.PublishedOn, err = parseDate(v)
			if err != nil {
				return nil, fmt.Errorf("%q is not a valid date", v)
			}
		case "updatedat":
		// do nothing
		default:
			return nil, fmt.Errorf("Unexpected key: %q", k)
		}
	}

	d, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	d, err = extractAdditionalMetadata(d, a)
	if err != nil {
		return nil, err
	}
	a.Body = d
	a.BodyHTML = msgToHTML(a.Body, a.Format)
	a.HTMLBody = template.HTML(a.BodyHTML)
	return a, nil
}

func readArticles() ([]*Article, []string, error) {
	timeStart := time.Now()
	walker := fs.Walk("articles")
	var res []*Article
	var dirs []string
	for walker.Step() {
		if walker.Err() != nil {
			fmt.Printf("readArticles: walker.Step() failed with %s\n", walker.Err())
			return nil, nil, walker.Err()
		}
		st := walker.Stat()
		path := walker.Path()
		if st.IsDir() {
			dirs = append(dirs, path)
			continue
		}
		name := filepath.Base(path)
		switch name {
		case "from-quicknotes.txt":
			continue
		}
		if name == "notes.txt" {
			err := readNotes(path)
			if err != nil {
				fmt.Printf("readWorkLog(%s) failed with %s\n", path, err)
				return nil, nil, err
			}
			continue
		}
		a, err := readArticle(path)
		if err != nil {
			fmt.Printf("readArticle() of %s failed with %s\n", path, err)
			return nil, nil, err
		}
		if a != nil {
			res = append(res, a)
		}
	}
	fmt.Printf("read %d articles in %s\n", len(res), time.Since(timeStart))
	return res, dirs, nil
}

// NewArticlesStore returns a store of articles
func NewArticlesStore() (*ArticlesStore, error) {
	articles, dirs, err := readArticles()
	if err != nil {
		return nil, err
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedOn.Before(articles[j].PublishedOn)
	})
	res := &ArticlesStore{articles: articles, dirsToWatch: dirs}
	res.idToArticle = make(map[string]*Article)
	for _, a := range articles {
		curr := res.idToArticle[a.ID]
		if curr == nil {
			res.idToArticle[a.ID] = a
			continue
		}
		log.Fatalf("2 articles with the same id %s\n%s\n%s\n", a.ID, curr.Path, a.Path)
	}
	return res, nil
}

// GetArticles returns all articles
func (s *ArticlesStore) GetArticles() []*Article {
	return s.articles
}

// GetArticleByID returns an article given its id
func (s *ArticlesStore) GetArticleByID(id string) *Article {
	return s.idToArticle[id]
}

// ArticlesCount returns number of articles
func (s *ArticlesStore) ArticlesCount() int {
	return len(s.articles)
}

// GetDirsToWatch returns directories to watch for chagnes
func (s *ArticlesStore) GetDirsToWatch() []string {
	return s.dirsToWatch
}

// TODO: this is simplistic but works for me, http://net.tutsplus.com/tutorials/other/8-regular-expressions-you-should-know/
// has more elaborate regex for extracting urls
var urlRx = regexp.MustCompile(`https?://[[:^space:]]+`)
var notURLEndChars = []byte(".),")

func notURLEndChar(c byte) bool {
	return -1 != bytes.IndexByte(notURLEndChars, c)
}

var disableUrlization = false

func strToHTML(s string) string {
	matches := urlRx.FindAllStringIndex(s, -1)
	if nil == matches || disableUrlization {
		s = template.HTMLEscapeString(s)
		s = strings.Replace(s, "\n", "<br>", -1)
		return "<p>" + s + "</p>"
	}

	urlMap := make(map[string]string)
	ns := ""
	prevEnd := 0
	for n, match := range matches {
		start, end := match[0], match[1]
		for end > start && notURLEndChar(s[end-1]) {
			end--
		}
		url := s[start:end]
		ns += s[prevEnd:start]

		// placeHolder is meant to be an unlikely string that doesn't exist in
		// the message, so that we can replace the string with it and then
		// revert the replacement. A more robust approach would be to remember
		// offsets
		placeHolder, ok := urlMap[url]
		if !ok {
			placeHolder = fmt.Sprintf("a;dfsl;a__lkasjdfh1234098;lajksdf_%d", n)
			urlMap[url] = placeHolder
		}
		ns += placeHolder
		prevEnd = end
	}
	ns += s[prevEnd:len(s)]

	ns = template.HTMLEscapeString(ns)
	for url, placeHolder := range urlMap {
		url = fmt.Sprintf(`<a href="%s" rel="nofollow">%s</a>`, url, url)
		ns = strings.Replace(ns, placeHolder, url, -1)
	}
	ns = strings.Replace(ns, "\n", "<br>", -1)
	return "<p>" + ns + "</p>"
}

func markdownToUnsafeHTML(text []byte) []byte {
	// Those are blackfriday.MarkdownCommon() extensions
	/*
		extensions := 0 |
			EXTENSION_NO_INTRA_EMPHASIS |
			EXTENSION_TABLES |
			EXTENSION_FENCED_CODE |
			EXTENSION_AUTOLINK |
			EXTENSION_STRIKETHROUGH |
			EXTENSION_SPACE_HEADERS |
			EXTENSION_HEADER_IDS |
			EXTENSION_BACKSLASH_LINE_BREAK |
			EXTENSION_DEFINITION_LISTS
	*/

	// https://github.com/shurcooL/github_flavored_markdown/blob/master/main.go#L82
	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

	commonHTMLFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	renderer := blackfriday.HtmlRenderer(commonHTMLFlags, "", "")
	opts := blackfriday.Options{Extensions: extensions}
	return blackfriday.MarkdownOptions(text, renderer, opts)
}

func markdownToHTML(s []byte) string {
	s, replacements := txtWithCodeParts(s)
	unsafe := markdownToUnsafeHTML(s)
	//unsafe := blackfriday.MarkdownCommon(s)
	policy := bluemonday.UGCPolicy()
	policy.AllowStyling()
	res := policy.SanitizeBytes(unsafe)
	for kstr, v := range replacements {
		k := []byte(kstr)
		res = bytes.Replace(res, k, v, -1)
	}
	return string(res)
}

func msgToHTML(msg []byte, format int) string {
	switch format {
	case formatHTML:
		return string(msg)
	case formatMarkdown:
		return markdownToHTML(msg)
	case formatText:
		return strToHTML(string(msg))
	}
	panic("unknown format")
}

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
	sha1 := u.Sha1HexOfBytes(jsData)
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
