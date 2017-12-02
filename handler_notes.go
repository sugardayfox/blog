package main

import (
	"net/http"
	"strings"

	"github.com/kjk/blog/pkg/notes"
)

type modelNotesForWeek struct {
	Notes         []*notes.Note
	TotalNotes    int
	TagCounts     []notes.TagWithCount
	WeekStartDay  string
	NextWeek      string
	PrevWeek      string
	AnalyticsCode string
}

// /dailynotes
func handleNotesIndex(w http.ResponseWriter, r *http.Request) {
	weekStart := notes.NotesWeekStarts[0]
	allNotes := notes.NotesWeekStartDayToNotes[weekStart]
	var nextWeek string
	if len(notes.NotesWeekStarts) > 1 {
		nextWeek = notes.NotesWeekStarts[1]
	}
	model := &modelNotesForWeek{
		Notes:         allNotes,
		TagCounts:     notes.NotesTagCounts,
		TotalNotes:    notes.TotalNotes,
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
	allNotes := notes.NotesWeekStartDayToNotes[weekStart]
	if len(allNotes) == 0 {
		serve404(w, r)
		return
	}
	var nextWeek, prevWeek string
	for idx, ws := range notes.NotesWeekStarts {
		if ws != weekStart {
			continue
		}
		if idx > 0 {
			prevWeek = notes.NotesWeekStarts[idx-1]
		}
		lastIdx := len(notes.NotesWeekStarts) - 1
		if idx+1 <= lastIdx {
			nextWeek = notes.NotesWeekStarts[idx+1]
		}
		break
	}
	model := &modelNotesForWeek{
		Notes:         allNotes,
		TagCounts:     notes.NotesTagCounts,
		WeekStartDay:  weekStart,
		NextWeek:      nextWeek,
		PrevWeek:      prevWeek,
		AnalyticsCode: analyticsCode,
	}
	serveTemplate(w, tmplNotesWeek, model)
}

// /worklog
func handleWorkLog(w http.ResponseWriter, r *http.Request) {
	// originally /dailynotes was under /worklog
	http.Redirect(w, r, "/dailynotes", http.StatusMovedPermanently)
}

// /dailynotes/note/${id}-${title}
func handleNotesNote(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	s := strings.TrimPrefix(uri, "/dailynotes/note/")
	parts := strings.SplitN(s, "-", 2)
	noteID := parts[0]
	aNote := notes.NotesIDToNote[noteID]
	if aNote == nil {
		serve404(w, r)
		return
	}

	weekStartTime := notes.CalcWeekStart(aNote.Day)
	weekStartDay := weekStartTime.Format("2006-01-02")
	model := struct {
		WeekStartDay  string
		Note          *notes.Note
		AnalyticsCode string
	}{
		WeekStartDay:  weekStartDay,
		Note:          aNote,
		AnalyticsCode: analyticsCode,
	}
	serveTemplate(w, tmplNotesNote, model)
}

// /dailynotes/tag/${tag} :
func handleNotesTag(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	tag := strings.TrimPrefix(uri, "/dailynotes/tag/")
	allNotes := notes.NotesTagToNotes[tag]

	if len(allNotes) == 0 {
		serve404(w, r)
		return
	}

	model := struct {
		Notes         []*notes.Note
		TagCounts     []notes.TagWithCount
		Tag           string
		AnalyticsCode string
	}{
		Notes:         allNotes,
		TagCounts:     notes.NotesTagCounts,
		Tag:           tag,
		AnalyticsCode: analyticsCode,
	}
	serveTemplate(w, tmplNotesTag, model)
}
