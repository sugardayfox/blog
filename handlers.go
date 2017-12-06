package main

import (
	"net/http"
	"time"
)

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}

	mux.HandleFunc("/", handleMainPage)

	mux.HandleFunc("/app/crashsubmit", handleCrashSubmit)
	mux.HandleFunc("/archives.html", handleArchives)
	mux.HandleFunc("/software", handleSoftware)
	mux.HandleFunc("/software/", handleSoftware)
	mux.HandleFunc("/extremeoptimizations/", handleExtremeOpt)
	mux.HandleFunc("/article/", handleArticle)
	mux.HandleFunc("/kb/", handleArticle)
	mux.HandleFunc("/blog/", handleArticle)
	mux.HandleFunc("/articles/", handleArticles)
	mux.HandleFunc("/book/", handleArticles)
	mux.HandleFunc("/tag/", handleTag)
	mux.HandleFunc("/static/", handleStatic)
	mux.HandleFunc("/dailynotes/week/", handleNotesWeek)
	mux.HandleFunc("/dailynotes/note/", handleNotesNote)
	mux.HandleFunc("/dailynotes", handleNotesIndex)
	mux.HandleFunc("/worklog", handleWorkLog)

	// not logged because not interesting for visitor analytics
	mux.HandleFunc("/ping", handlePing)
	mux.HandleFunc("/css/", handleCSS)
	mux.HandleFunc("/js/", handleJs)
	mux.HandleFunc("/gfx/", handleGfx)

	mux.HandleFunc("/djs/", handleDjs)

	// websocket is only for dev mode, used for refreshing the pages if
	// they change on disk
	if !flgProduction {
		mux.HandleFunc("/ws", serveWs)
	}

	// https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv
}
