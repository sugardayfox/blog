package main

import (
	"net/http"
	"time"
)

// https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}

	mux.HandleFunc("/", handleMainPage)
	mux.HandleFunc("/favicon.ico", handleFavicon)
	mux.HandleFunc("/robots.txt", handleRobotsTxt)
	mux.HandleFunc("/contactme.html", handleContactme)
	mux.HandleFunc("/oauthtwittercb", handleOauthTwitterCallback)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/logout", handleLogout)

	mux.HandleFunc("/app/crashsubmit", handleCrashSubmit)
	mux.HandleFunc("/app/crashes", handleCrashes)
	mux.HandleFunc("/app/crashesrss", handleCrashesRss)
	mux.HandleFunc("/app/crashshow", handleCrashShow)
	// TODO: I stopped pointing people to FeedBurner feed on 2013-05-22
	// At some point I should delete /feedburner.xml, which is a source data
	// for FeedBurner
	mux.HandleFunc("/feedburner.xml", handleAtom)
	mux.HandleFunc("/atom.xml", handleAtom)
	mux.HandleFunc("/atom-all.xml", handleAtomAll)
	mux.HandleFunc("/archives.html", handleArchives)
	mux.HandleFunc("/software", handleSoftware)
	mux.HandleFunc("/software/", handleSoftware)
	mux.HandleFunc("/extremeoptimizations/", handleExtremeOpt)
	mux.HandleFunc("/article/", handleArticle)
	mux.HandleFunc("/kb/", handleArticle)
	mux.HandleFunc("/blog/", handleArticle)
	mux.HandleFunc("/forum_sumatra/", forumRedirect)
	mux.HandleFunc("/articles/", handleArticles)
	mux.HandleFunc("/tag/", handleTag)
	mux.HandleFunc("/static/", handleStatic)
	mux.HandleFunc("/css/", handleCss)
	mux.HandleFunc("/js/", handleJs)
	mux.HandleFunc("/gfx/", handleGfx)
	mux.HandleFunc("/djs/", handleDjs)
	if !inProduction {
		mux.HandleFunc("/ws", serveWs)
	}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		// TODO: 1.8 only
		// IdleTimeout:  120 * time.Second,
		Handler: mux,
	}
	// TODO: track connections and their state
	return srv
}
