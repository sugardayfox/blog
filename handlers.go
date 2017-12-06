package main

import (
	"net/http"
	"time"
)

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}

	mux.HandleFunc("/app/crashsubmit", handleCrashSubmit)
	mux.HandleFunc("/extremeoptimizations/", handleExtremeOpt)
	mux.HandleFunc("/articles/", handleArticles)
	mux.HandleFunc("/book/", handleArticles)
	mux.HandleFunc("/static/", handleStatic)

	mux.HandleFunc("/djs/", handleDjs)

	// https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv
}
