package main

import (
	"bufio"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/kjk/u"
	"github.com/mvdan/xurls"
)

const (
	partSeparator = "---"
)

var (
	workLogDays        []*workLogDay
	workLogTagsToParts map[string][]*workLogPart
)

type workLogPart struct {
	Day      time.Time
	BodyHTML string
	Tags     []string
}

type workLogDay struct {
	Day    time.Time
	DayStr string
	Parts  []*workLogPart
}

// RemoveDuplicateStrings removes duplicate strings from a
// Is optimized for the case of no duplicates
func RemoveDuplicateStrings(a []string) []string {
	sort.Strings(a)
	hasDups := false
	for i := 1; i < len(a); i++ {
		if a[i-1] == a[i] {
			hasDups = true
			break
		}
	}
	if !hasDups {
		return a
	}
	var res []string
	m := make(map[string]struct{})
	for _, s := range a {
		if _, ok := m[s]; !ok {
			m[s] = struct{}{}
			res = append(res, s)
		}
	}
	return res
}

// a line is just a #hashtag if it has only one word and starts with #
func isJustHashtag(s string) bool {
	if !strings.HasPrefix(s, "#") {
		return false
	}
	// Maybe: consider more characters other than ' ' as #hashtag delimiter
	return !strings.Contains(s, " ")
}

func lastLineEmpty(lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	lastIdx := len(lines) - 1
	line := lines[lastIdx]
	return len(line) == 0
}

func removeLastLine(lines []string) []string {
	lastIdx := len(lines) - 1
	return lines[:lastIdx]
}

func buildBodyFromLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	// remove empty lines from beginning
	for len(lines[0]) == 0 {
		lines = lines[1:]
	}

	// remove empty lines from end
	for lastLineEmpty(lines) {
		lines = removeLastLine(lines)
	}
	if len(lines) == 0 {
		return ""
	}

	// collapse multiple empty lines into single empty line
	// and remove lines that are just #hashtags
	currWrite := 1
	for i := 1; i < len(lines); i++ {
		prev := lines[currWrite-1]
		curr := lines[i]
		if len(prev) == 0 && len(curr) == 0 {
			// skips the current line because we don't advance currWrite
			continue
		}
		if isJustHashtag(curr) {
			// skip just hashtags
			continue
		}

		if i != currWrite {
			lines[currWrite] = lines[i]
		}
		currWrite++
	}
	lines = lines[:currWrite]
	return strings.Join(lines, "\n")
}

// tags start with #
func extractTagsFromString(txt string) []string {
	var res []string
	parts := strings.Split(txt, " ")
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "#") {
			res = append(res, s[1:])
		}
	}
	return res
}

func extractTagsFromLines(lines []string) []string {
	var res []string
	for _, line := range lines {
		tags := extractTagsFromString(line)
		res = append(res, tags...)
	}
	return res
}

// there are no guarantees in live, but this should be pretty unique string
func genRandomString() string {
	var a [20]byte
	buf := a[:]
	_, err := rand.Read(buf)
	if err == nil {
		return string(buf)
	}
	return fmt.Sprintf("__--##%d##--__", rand.Int63())
}

func workLogPostToHTML(s string) string {
	urls := xurls.Relaxed.FindAllString(s, -1)
	urls = RemoveDuplicateStrings(urls)

	// sort by length, longest first, so that we correctly convert
	// urls to hrefs when there are 2 urls like http://foo.com
	// and http://foo.com/longer
	sort.Slice(urls, func(i, j int) bool {
		return len(urls[i]) > len(urls[j])
	})
	// this is a two-step url -> random_unique_string,
	// random_unique_string -> url replacement to prevent
	// double-escaping if we have 2 urls like: foo.bar.com and bar.com
	urlToUnique := make(map[string]string)

	for _, url := range urls {
		unique := genRandomString()
		urlToUnique[url] = unique
		s = strings.Replace(s, url, unique, -1)
	}

	for _, url := range urls {
		replacement := fmt.Sprintf(`<a href="%s">%s</a>`, url, url)
		unique := urlToUnique[url]
		s = strings.Replace(s, unique, replacement, -1)
	}

	s, _ = sanitize.HTMLAllowing(s, []string{"a"})
	return s
}

func newWorkLogPart(lines []string) *workLogPart {
	tags := extractTagsFromLines(lines)
	s := buildBodyFromLines(lines)
	body := workLogPostToHTML(s)
	return &workLogPart{
		Tags:     tags,
		BodyHTML: body,
	}
}

func workLogLinesToParts(lines []string) []*workLogPart {
	// parts are separated by "---" line
	var res []*workLogPart
	var curr []string
	for _, line := range lines {
		if line == partSeparator {
			if len(curr) > 0 {
				part := newWorkLogPart(curr)
				res = append(res, part)
			}
			curr = nil
		} else {
			curr = append(curr, line)
		}
	}
	if len(curr) > 0 {
		part := newWorkLogPart(curr)
		res = append(res, part)
	}
	return res
}

func readWorkLog(path string) error {
	workLogTagsToParts = make(map[string][]*workLogPart)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var posts []*workLogDay
	var curr *workLogDay
	var lines []string

	for scanner.Scan() {
		s := strings.TrimRight(scanner.Text(), "\n\r\t ")
		day, err := time.Parse("2006-01-02", s)

		if err != nil {
			// first line must be a valid new day
			u.PanicIf(curr == nil)
			lines = append(lines, s)
			continue
		}

		// this is a new day
		if curr != nil {
			curr.Parts = workLogLinesToParts(lines)
			posts = append(posts, curr)
		}
		curr = &workLogDay{
			Day:    day,
			DayStr: s,
		}
		lines = nil
	}
	curr.Parts = workLogLinesToParts(lines)
	workLogDays = append(posts, curr)

	// verify they are in chronological order
	for i := 1; i < len(workLogDays); i++ {
		post := workLogDays[i-1]
		postPrev := workLogDays[i]
		diff := post.Day.Sub(postPrev.Day)
		if diff < 0 {
			return fmt.Errorf("Post '%s' should be later than '%s'", post.DayStr, postPrev.DayStr)
		}
	}
	// update date on parts
	for _, post := range workLogDays {
		for _, part := range post.Parts {
			part.Day = post.Day
			for _, tag := range part.Tags {
				a := workLogTagsToParts[tag]
				a = append(a, part)
				workLogTagsToParts[tag] = a
			}
		}
	}
	fmt.Printf("Read %d daily logs\n", len(workLogDays))
	// TODO: build weekly index (week starting on Monday)
	// TODO: serve weekly posts, anchored at /worklog/, /worklog/{first-day-of-week}
	return scanner.Err()
}

type modelWorkLogPost struct {
	DayStr   string
	HTMLBody template.HTML
	Tags     []string
}

type modelWorkLogDay struct {
	DayStr string
	Posts  []*modelWorkLogPost
}

type modelWorkLogIndex struct {
	Days          []*modelWorkLogDay
	StartDay      string
	NextWeek      string
	PrevWeek      string
	AnalyticsCode string
}

func makeModelWorkLogPost(post *workLogPart) *modelWorkLogPost {
	dayStr := post.Day.Format("2006-01-02 Mon")
	return &modelWorkLogPost{
		DayStr:   dayStr,
		HTMLBody: template.HTML(post.BodyHTML),
		Tags:     post.Tags,
	}
}

// /worklog
func handleWorkLogIndex(w http.ResponseWriter, r *http.Request) {
	var days []*modelWorkLogDay
	for _, day := range workLogDays {
		var posts []*modelWorkLogPost
		for _, post := range day.Parts {
			posts = append(posts, makeModelWorkLogPost(post))
		}
		day2 := &modelWorkLogDay{
			Posts: posts,
		}
		days = append(days, day2)
	}

	// for index page we don't set Start Day
	model := &modelWorkLogIndex{
		Days:          days,
		AnalyticsCode: analyticsCode,
		// TODO: NextWeek, PrevWeek
	}
	execTemplate(w, tmplWorkLogWeek, model)
}

// /worklog/week/${day} : week starting with a given day
func handleWorkLogWeek(w http.ResponseWriter, r *http.Request) {
	serve404(w, r)
}

type modelWorkLogTag struct {
	Posts         []*modelWorkLogPost
	Tag           string
	AnalyticsCode string
}

// /worklog/tag/${tag} :
func handleWorkLogTag(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	tag := strings.TrimPrefix(uri, "/worklog/tag/")
	posts := workLogTagsToParts[tag]

	if len(posts) == 0 {
		serve404(w, r)
		return
	}
	var postsModel []*modelWorkLogPost
	for _, post := range posts {
		postsModel = append(postsModel, makeModelWorkLogPost(post))
	}
	model := &modelWorkLogTag{
		Posts:         postsModel,
		Tag:           tag,
		AnalyticsCode: analyticsCode,
	}
	execTemplate(w, tmplWorkLogTag, model)
}
