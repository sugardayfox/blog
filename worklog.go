package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/kjk/u"
	"github.com/mvdan/xurls"
)

const (
	postSeparator = "---"
)

var (
	workLogDays        []*workLogDay
	workLogTagsToPosts map[string][]*workLogPost

	workLogWeekStartDayToPosts map[string][]*workLogPost
	workLogWeekStarts          []string
)

type workLogPost struct {
	Day    time.Time
	DayStr string // in format 2006-01-02
	// in format 2006-01-02-${idx}. This is an index within workLogDay.Posts
	// which is not ideal because it changes if I delete a post or re-arrange
	// them, but that's rare. The alternative would be to auto-generate
	// unique ids, e.g. parsing would add missing data and re-save
	ID       string
	HTMLBody string
	Tags     []string
}

type workLogDay struct {
	Day    time.Time
	DayStr string
	Posts  []*workLogPost
}

type modelWorkLogWeek struct {
	Posts         []*workLogPost
	WeekStartDay  string
	NextWeek      string
	PrevWeek      string
	AnalyticsCode string
}

type modelWorkLogTag struct {
	Posts         []*workLogPost
	Tag           string
	AnalyticsCode string
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

func findWordEnd(s string, start int) int {
	for i := start; i < len(s); i++ {
		c := s[i]
		if c == ' ' {
			return i + 1
		}
	}
	return -1
}

// TODO: must not remove spaces from start
func collapseMultipleSpaces(s string) string {
	for {
		s2 := strings.Replace(s, "  ", " ", -1)
		if s2 == s {
			return s

		}
		s = s2
	}
}

// remove hashtags from start and end
func removeHashtags(s string) string {
	// remove hashtags from start
	for strings.HasPrefix(s, "#") {
		idx := findWordEnd(s, 0)
		if idx == -1 {
			return ""
		}
		s = s[idx:]
		s = strings.TrimLeft(s, " ")
	}

	// remove hashtags from end
	for {
		idx := strings.LastIndex(s, "#")
		if idx == -1 {
			return s
		}
		if -1 != findWordEnd(s, idx) {
			return s
		}
		s = strings.TrimRight(s[:idx], " ")
	}
}

func buildBodyFromLines(lines []string) string {
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

		lines[currWrite] = curr
		currWrite++
	}
	lines = lines[:currWrite]
	for idx, line := range lines {
		lines[idx] = removeHashtags(line)
	}

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
	urls = u.RemoveDuplicateStrings(urls)

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

func newWorkLogPost(lines []string) *workLogPost {
	tags := extractTagsFromLines(lines)
	s := buildBodyFromLines(lines)
	body := workLogPostToHTML(s)
	return &workLogPost{
		Tags:     tags,
		HTMLBody: body,
	}
}

func workLogLinesToPosts(lines []string) []*workLogPost {
	// parts are separated by "---" line
	var res []*workLogPost
	var curr []string
	for _, line := range lines {
		if line == postSeparator {
			if len(curr) > 0 {
				part := newWorkLogPost(curr)
				res = append(res, part)
			}
			curr = nil
		} else {
			curr = append(curr, line)
		}
	}
	if len(curr) > 0 {
		part := newWorkLogPost(curr)
		res = append(res, part)
	}
	return res
}

func readWorkLog(path string) error {
	workLogTagsToPosts = make(map[string][]*workLogPost)
	workLogWeekStartDayToPosts = make(map[string][]*workLogPost)
	workLogWeekStarts = nil
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
			curr.Posts = workLogLinesToPosts(lines)
			posts = append(posts, curr)
		}
		curr = &workLogDay{
			Day:    day,
			DayStr: s,
		}
		lines = nil
	}
	curr.Posts = workLogLinesToPosts(lines)
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

	// update date and id on posts
	for _, day := range workLogDays {
		weekStartTime := calcWeekStart(day.Day)
		weekStartDay := weekStartTime.Format("2006-01-02")
		for idx, post := range day.Posts {
			post.Day = day.Day
			post.DayStr = day.Day.Format("2006-01-02")
			post.ID = fmt.Sprintf("%s-%d", post.DayStr, idx)
			for _, tag := range post.Tags {
				a := workLogTagsToPosts[tag]
				a = append(a, post)
				workLogTagsToPosts[tag] = a
			}
			a := workLogWeekStartDayToPosts[weekStartDay]
			a = append(a, post)
			workLogWeekStartDayToPosts[weekStartDay] = a
		}
	}
	for day := range workLogWeekStartDayToPosts {
		workLogWeekStarts = append(workLogWeekStarts, day)
	}
	sort.Strings(workLogWeekStarts)
	fmt.Printf("Read %d daily logs\n", len(workLogDays))
	fmt.Printf("workLogWeekStarts: %v\n", workLogWeekStarts)
	return scanner.Err()
}

// given time, return time on start of week (monday)
func calcWeekStart(t time.Time) time.Time {
	// wd is 1 to 7
	wd := t.Weekday()
	dayOffset := time.Duration((wd - 1)) * time.Hour * -24
	return t.Add(dayOffset)
}

// /worklog
func handleWorkLogIndex(w http.ResponseWriter, r *http.Request) {
	weekStart := workLogWeekStarts[0]
	posts := workLogWeekStartDayToPosts[weekStart]
	var nextWeek string
	if len(workLogWeekStarts) > 1 {
		nextWeek = workLogWeekStarts[1]
	}
	model := &modelWorkLogWeek{
		Posts:         posts,
		WeekStartDay:  weekStart,
		AnalyticsCode: analyticsCode,
		NextWeek:      nextWeek,
	}
	execTemplate(w, tmplWorkLogWeek, model)
}

// /worklog/week/${day} : week starting with a given day
func handleWorkLogWeek(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	weekStart := strings.TrimPrefix(uri, "/worklog/week/")
	posts := workLogWeekStartDayToPosts[weekStart]
	if len(posts) == 0 {
		serve404(w, r)
		return
	}
	var nextWeek, prevWeek string
	for idx, ws := range workLogWeekStarts {
		if ws != weekStart {
			continue
		}
		if idx > 0 {
			prevWeek = workLogWeekStarts[idx-1]
		}
		lastIdx := len(workLogWeekStarts) - 1
		if idx+1 <= lastIdx {
			nextWeek = workLogWeekStarts[idx+1]
		}
		break
	}
	model := &modelWorkLogWeek{
		Posts:         posts,
		WeekStartDay:  weekStart,
		NextWeek:      nextWeek,
		PrevWeek:      prevWeek,
		AnalyticsCode: analyticsCode,
	}
	execTemplate(w, tmplWorkLogWeek, model)
}

func findWorkLogDay(dayStr string) *workLogDay {
	for _, d := range workLogDays {
		if dayStr == d.DayStr {
			return d
		}
	}
	return nil
}

// /worklog/post/${day}-${idx}
func handleWorkLogPost(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	postID := strings.TrimPrefix(uri, "/worklog/post/")
	// expecting sth. like: 2006-01-02-1
	parts := strings.Split(postID, "-")
	if len(parts) != 4 {
		serve404(w, r)
		return
	}
	idx, err := strconv.Atoi(parts[3])
	if err != nil || idx < 0 {
		serve404(w, r)
		return
	}

	dateStr := strings.Join(parts[:3], "-")
	day := findWorkLogDay(dateStr)
	if day == nil {
		serve404(w, r)
		return
	}

	if idx >= len(day.Posts) {
		serve404(w, r)
		return
	}

	post := day.Posts[idx]
	weekStartTime := calcWeekStart(day.Day)
	weekStartDay := weekStartTime.Format("2006-01-02")
	model := struct {
		WeekStartDay  string
		Post          *workLogPost
		AnalyticsCode string
	}{
		WeekStartDay:  weekStartDay,
		Post:          post,
		AnalyticsCode: analyticsCode,
	}
	execTemplate(w, tmplWorkLogPost, model)
}

// /worklog/tag/${tag} :
func handleWorkLogTag(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	tag := strings.TrimPrefix(uri, "/worklog/tag/")
	posts := workLogTagsToPosts[tag]

	if len(posts) == 0 {
		serve404(w, r)
		return
	}
	model := &modelWorkLogTag{
		Posts:         posts,
		Tag:           tag,
		AnalyticsCode: analyticsCode,
	}
	execTemplate(w, tmplWorkLogTag, model)
}
