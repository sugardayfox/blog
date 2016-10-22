package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kr/fs"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type Article struct {
	Id          int
	PublishedOn time.Time
	Title       string
	Tags        []string
	Format      int
	Path        string
	Body        []byte
	BodyHtml    string
}

const (
	FormatHtml     = 0
	FormatMarkdown = 2
	FormatText     = 3

	FormatFirst   = 0
	FormatLast    = 3
	FormatUnknown = -1
)

type ArticlesByTime []*Article

func (s ArticlesByTime) Len() int {
	return len(s)
}

func (s ArticlesByTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ArticlesByTime) Less(i, j int) bool {
	return s[i].PublishedOn.Before(s[j].PublishedOn)
}

// same format as Format* constants
var formatNames = []string{"Html", "Textile", "Markdown", "Text"}

func validFormat(format int) bool {
	return format >= FormatFirst && format <= FormatLast
}

func remSep(s string) string {
	return strings.Replace(s, "|", "", -1)
}

func urlForTag(tag string) string {
	// TODO: url-quote the first tag
	return fmt.Sprintf(`<a href="/tag/%s" class="taglink">%s</a>`, tag, tag)
}

func FormatNameToId(name string) int {
	for i, formatName := range formatNames {
		if strings.EqualFold(name, formatName) {
			return i
		}
	}
	return FormatUnknown
}

type Store struct {
	articles    []*Article
	idToArticle map[int]*Article
	dirsToWatch []string
}

func isSepLine(s string) bool {
	return strings.HasPrefix(s, "-----")
}

func parseTags(s string) []string {
	tags := strings.Split(s, ",")
	for i, tag := range tags {
		tag = strings.TrimSpace(tag)
		tag = strings.ToLower(tag)
		tags[i] = tag
	}
	return tags
}

func parseFormat(s string) int {
	s = strings.ToLower(s)
	switch s {
	case "html":
		return FormatHtml
	case "markdown", "md":
		return FormatMarkdown
	case "text":
		return FormatText
	default:
		return FormatUnknown
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
			return nil, fmt.Errorf("Unexpected line: %q\n", l)
		}
		k := strings.ToLower(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "deleted", "draft":
			if inProduction {
				return nil, nil
			}
		case "id":
			id, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("%q is not a valid id (not a number)", v)
			}
			a.Id = id
			a.Path = path
		case "title":
			a.Title = v
		case "tags":
			a.Tags = parseTags(v)
		case "format":
			f := parseFormat(v)
			if f == FormatUnknown {
				return nil, fmt.Errorf("%q is not a valid format", v)
			}
			a.Format = f
		case "date":
			a.PublishedOn, err = parseDate(v)
			if err != nil {
				return nil, fmt.Errorf("%q is not a valid date", v)
			}
		default:
			return nil, fmt.Errorf("Unexpected key: %q\n", k)
		}
	}
	a.Body, err = ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func readArticles() ([]*Article, []string, error) {
	timeStart := time.Now()
	walker := fs.Walk("blog_posts")
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

func NewStore() (*Store, error) {
	articles, dirs, err := readArticles()
	if err != nil {
		return nil, err
	}
	sort.Sort(ArticlesByTime(articles))
	res := &Store{articles: articles, dirsToWatch: dirs}
	res.idToArticle = make(map[int]*Article)
	for _, a := range articles {
		curr := res.idToArticle[a.Id]
		if curr == nil {
			res.idToArticle[a.Id] = a
			continue
		}
		log.Fatalf("2 articles with the same id %d\n%s\n%s\n", a.Id, curr.Path, a.Path)
	}
	return res, nil
}

func (s *Store) GetArticles() []*Article {
	return s.articles
}

func (s *Store) GetArticleById(id int) *Article {
	//fmt.Printf("GetArticleById: %d\n", id)
	for _, a := range s.articles {
		if a.Id == id {
			return a
		}
	}
	return nil
}

func (s *Store) ArticlesCount() int {
	return len(s.articles)
}

func (a *Article) Permalink() string {
	return "article/" + ShortenId(a.Id) + "/" + Urlify(a.Title) + ".html"
}

func (a *Article) TagsDisplay() template.HTML {
	arr := make([]string, 0)
	for _, tag := range a.Tags {
		arr = append(arr, urlForTag(tag))
	}
	s := strings.Join(arr, ", ")
	return template.HTML(s)
}

func (a *Article) GetHtmlStr() string {
	if a.BodyHtml == "" {
		a.BodyHtml = msgToHTML(a.Body, a.Format)
	}
	return a.BodyHtml
}

func (s *Store) GetDirsToWatch() []string {
	return s.dirsToWatch
}

// TODO: this is simplistic but works for me, http://net.tutsplus.com/tutorials/other/8-regular-expressions-you-should-know/
// has more elaborate regex for extracting urls
var urlRx = regexp.MustCompile(`https?://[[:^space:]]+`)
var notUrlEndChars = []byte(".),")

func notUrlEndChar(c byte) bool {
	return -1 != bytes.IndexByte(notUrlEndChars, c)
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
		for end > start && notUrlEndChar(s[end-1]) {
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

func markdown(s []byte) string {
	s, replacements := txtWithCodeParts(s)
	unsafe := blackfriday.MarkdownCommon(s)
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
	case FormatHtml:
		return string(msg)
	case FormatMarkdown:
		return markdown(msg)
	case FormatText:
		return strToHTML(string(msg))
	}
	panic("unknown format")
}
