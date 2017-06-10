package main

import (
	"bufio"
	"encoding/hex"
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
	"github.com/sourcegraph/syntaxhighlight"
)

const (
	noteSeparator  = "---"
	codeBlockStart = "```"
)

var (
	notesDays        []*notesForDay
	notesTagsToNotes map[string][]*note

	notesWeekStartDayToNotes map[string][]*note
	notesWeekStarts          []string
	nTotalNotes              int
)

type note struct {
	Day    time.Time
	DayStr string // in format 2006-01-02
	// in format 2006-01-02-${idx}. This is an index within notesForDay.Notes
	// which is not ideal because it changes if I delete a post or re-arrange
	// them, but that's rare. The alternative would be to auto-generate
	// unique ids, e.g. parsing would add missing data and re-save
	ID       string
	HTMLBody string
	Tags     []string
}

type notesForDay struct {
	Day    time.Time
	DayStr string
	Notes  []*note
}

type modelNotesForWeek struct {
	Notes         []*note
	TotalNotes    int
	WeekStartDay  string
	NextWeek      string
	PrevWeek      string
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
// TODO: maybe only at the beginning/end?
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
	_, err := rand.Read(a[:])
	if err == nil {
		return hex.EncodeToString(a[:])
	}
	return fmt.Sprintf("__--##%d##--__", rand.Int63())
}

func noteToHTML(s string) string {
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
	urlToAnchor := make(map[string]string)

	for _, url := range urls {
		anchor := genRandomString()
		urlToAnchor[url] = anchor
		s = strings.Replace(s, url, anchor, -1)
	}

	for _, url := range urls {
		replacement := fmt.Sprintf(`<a href="%s">%s</a>`, url, url)
		anchor := urlToAnchor[url]
		s = strings.Replace(s, anchor, replacement, -1)
	}

	s, _ = sanitize.HTMLAllowing(s, []string{"a"})
	return s
}

// returns new lines and a mapping of string => html as flattened string array
func extractCodeSnippets(lines []string) ([]string, []string) {
	var resLines []string
	var anchors []string
	codeLineStart := -1
	for i, s := range lines {
		isCodeLine := strings.HasPrefix(s, codeBlockStart)
		if isCodeLine {
			if codeLineStart == -1 {
				// this is a beginning of new code block
				codeLineStart = i
			} else {
				// end of the code block
				//lang := strings.TrimPrefix(lines[codeLineStart], codeBlockStart)
				codeLines := lines[codeLineStart+1 : i]
				codeLineStart = -1
				code := strings.Join(codeLines, "\n")
				codeHTML, err := syntaxhighlight.AsHTML([]byte(code))
				u.PanicIfErr(err)
				anchor := genRandomString()
				resLines = append(resLines, anchor)
				anchors = append(anchors, anchor, string(codeHTML))
			}
		} else {
			if codeLineStart == -1 {
				resLines = append(resLines, s)
			}
		}
	}
	// TODO: could append unclosed lines
	u.PanicIf(codeLineStart != -1)

	return resLines, anchors
}

func newNote(lines []string) *note {
	nTotalNotes++
	tags := extractTagsFromLines(lines)
	lines, codeReplacements := extractCodeSnippets(lines)
	s := buildBodyFromLines(lines)
	body := noteToHTML(s)
	n := len(codeReplacements) / 2
	for i := 0; i < n; i++ {
		anchor := codeReplacements[i*2]
		codeHTML := `<pre class="note-code">` + codeReplacements[i*2+1] + `</pre>`
		body = strings.Replace(body, anchor, codeHTML, -1)
	}
	return &note{
		Tags:     tags,
		HTMLBody: body,
	}
}

func linesToNotes(lines []string) []*note {
	// parts are separated by "---" line
	var res []*note
	var curr []string
	for _, line := range lines {
		if line == noteSeparator {
			if len(curr) > 0 {
				part := newNote(curr)
				res = append(res, part)
			}
			curr = nil
		} else {
			curr = append(curr, line)
		}
	}
	if len(curr) > 0 {
		part := newNote(curr)
		res = append(res, part)
	}
	return res
}

func readNotes(path string) error {
	notesTagsToNotes = make(map[string][]*note)
	notesWeekStartDayToNotes = make(map[string][]*note)
	notesWeekStarts = nil
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var posts []*notesForDay
	var curr *notesForDay
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
			curr.Notes = linesToNotes(lines)
			posts = append(posts, curr)
		}
		curr = &notesForDay{
			Day:    day,
			DayStr: s,
		}
		lines = nil
	}
	curr.Notes = linesToNotes(lines)
	notesDays = append(posts, curr)

	// verify they are in chronological order
	for i := 1; i < len(notesDays); i++ {
		post := notesDays[i-1]
		postPrev := notesDays[i]
		diff := post.Day.Sub(postPrev.Day)
		if diff < 0 {
			return fmt.Errorf("Post '%s' should be later than '%s'", post.DayStr, postPrev.DayStr)
		}
	}

	// update date and id on posts
	for _, day := range notesDays {
		weekStartTime := calcWeekStart(day.Day)
		weekStartDay := weekStartTime.Format("2006-01-02")
		for idx, post := range day.Notes {
			post.Day = day.Day
			post.DayStr = day.Day.Format("2006-01-02")
			post.ID = fmt.Sprintf("%s-%d", post.DayStr, idx)
			for _, tag := range post.Tags {
				a := notesTagsToNotes[tag]
				a = append(a, post)
				notesTagsToNotes[tag] = a
			}
			a := notesWeekStartDayToNotes[weekStartDay]
			a = append(a, post)
			notesWeekStartDayToNotes[weekStartDay] = a
		}
	}
	for day := range notesWeekStartDayToNotes {
		notesWeekStarts = append(notesWeekStarts, day)
	}
	sort.Strings(notesWeekStarts)
	fmt.Printf("Read %d daily logs\n", len(notesDays))
	fmt.Printf("notesWeekStarts: %v\n", notesWeekStarts)
	return scanner.Err()
}

// given time, return time on start of week (monday)
func calcWeekStart(t time.Time) time.Time {
	// wd is 1 to 7
	wd := t.Weekday()
	dayOffset := time.Duration((wd - 1)) * time.Hour * -24
	return t.Add(dayOffset)
}

// /dailynotes
func handleNotesIndex(w http.ResponseWriter, r *http.Request) {
	weekStart := notesWeekStarts[0]
	posts := notesWeekStartDayToNotes[weekStart]
	var nextWeek string
	if len(notesWeekStarts) > 1 {
		nextWeek = notesWeekStarts[1]
	}
	model := &modelNotesForWeek{
		Notes:         posts,
		TotalNotes:    nTotalNotes,
		WeekStartDay:  weekStart,
		AnalyticsCode: analyticsCode,
		NextWeek:      nextWeek,
	}
	serveTemplate(w, tmplNotesWeek, model)
}

// /dailynotes/week/${day} : week starting with a given day
func handleNotesWeek(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	weekStart := strings.TrimPrefix(uri, "/dailynotes/week/")
	posts := notesWeekStartDayToNotes[weekStart]
	if len(posts) == 0 {
		serve404(w, r)
		return
	}
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
		Notes:         posts,
		WeekStartDay:  weekStart,
		NextWeek:      nextWeek,
		PrevWeek:      prevWeek,
		AnalyticsCode: analyticsCode,
	}
	serveTemplate(w, tmplNotesWeek, model)
}

func findNotesForDay(dayStr string) *notesForDay {
	for _, d := range notesDays {
		if dayStr == d.DayStr {
			return d
		}
	}
	return nil
}

// /dailynotes/note/${day}-${idx}
func handleNotesNote(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	postID := strings.TrimPrefix(uri, "/dailynotes/note/")
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
	day := findNotesForDay(dateStr)
	if day == nil {
		serve404(w, r)
		return
	}

	if idx >= len(day.Notes) {
		serve404(w, r)
		return
	}

	oneNote := day.Notes[idx]
	weekStartTime := calcWeekStart(day.Day)
	weekStartDay := weekStartTime.Format("2006-01-02")
	model := struct {
		WeekStartDay  string
		Note          *note
		AnalyticsCode string
	}{
		WeekStartDay:  weekStartDay,
		Note:          oneNote,
		AnalyticsCode: analyticsCode,
	}
	serveTemplate(w, tmplNotesNote, model)
}

// /dailynotes/tag/${tag} :
func handleNotesTag(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	tag := strings.TrimPrefix(uri, "/dailynotes/tag/")
	posts := notesTagsToNotes[tag]

	if len(posts) == 0 {
		serve404(w, r)
		return
	}

	model := struct {
		Notes         []*note
		Tag           string
		AnalyticsCode string
	}{
		Notes:         posts,
		Tag:           tag,
		AnalyticsCode: analyticsCode,
	}
	serveTemplate(w, tmplNotesTag, model)
}
