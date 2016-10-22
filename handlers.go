package main

import "net/http"

func InitHttpHandlers() {
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/favicon.ico", handleFavicon)
	http.HandleFunc("/robots.txt", handleRobotsTxt)
	http.HandleFunc("/contactme.html", handleContactme)
	http.HandleFunc("/oauthtwittercb", handleOauthTwitterCallback)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)

	http.HandleFunc("/app/crashsubmit", handleCrashSubmit)
	http.HandleFunc("/app/crashes", handleCrashes)
	http.HandleFunc("/app/crashesrss", handleCrashesRss)
	http.HandleFunc("/app/crashshow", handleCrashShow)
	// TODO: I stopped pointing people to FeedBurner feed on 2013-05-22
	// At some point I should delete /feedburner.xml, which is a source data
	// for FeedBurner
	http.HandleFunc("/feedburner.xml", handleAtom)
	http.HandleFunc("/atom.xml", handleAtom)
	http.HandleFunc("/atom-all.xml", handleAtomAll)
	http.HandleFunc("/archives.html", handleArchives)
	http.HandleFunc("/software", handleSoftware)
	http.HandleFunc("/software/", handleSoftware)
	http.HandleFunc("/extremeoptimizations/", handleExtremeOpt)
	http.HandleFunc("/article/", handleArticle)
	http.HandleFunc("/kb/", handleArticle)
	http.HandleFunc("/blog/", handleArticle)
	http.HandleFunc("/forum_sumatra/", forumRedirect)
	http.HandleFunc("/articles/", handleArticles)
	http.HandleFunc("/tag/", handleTag)
	http.HandleFunc("/static/", handleStatic)
	http.HandleFunc("/css/", handleCss)
	http.HandleFunc("/js/", handleJs)
	http.HandleFunc("/gfx/", handleGfx)
	http.HandleFunc("/markitup/", handleMarkitup)
	http.HandleFunc("/djs/", handleDjs)
	if !inProduction {
		http.HandleFunc("/ws", serveWs)
	}
}
