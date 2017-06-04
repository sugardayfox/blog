package main

import (
	"net/http"
	"path/filepath"
	"strings"
)

func getWwwDir() string {
	// when running locally
	d := filepath.Join("..", "www")
	if pathExists(d) {
		return d
	}
	// when running on a server
	d = "www"
	if pathExists(d) {
		return d
	}
	logger.Errorf("getWwwDir(): %q dir doesn't exist", d)
	return ""
}

func getAppEngineTmplDir() string {
	// when running locally
	d := filepath.Join("..", "tmpl")
	if pathExists(d) {
		return d
	}
	// when running on a server
	d = "appengtmpl"
	if pathExists(d) {
		return d
	}
	logger.Errorf("getAppEngineTmplDir(): %q dir doesn't exist", d)
	return ""
}

func getCSSDir() string {
	return filepath.Join(getWwwDir(), "css")
}

func getJsDir() string {
	return filepath.Join(getWwwDir(), "js")
}

func getGfxDir() string {
	return filepath.Join(getWwwDir(), "gfx")
}

func getStaticDir() string {
	return filepath.Join(getWwwDir(), "static")
}

func getSoftwareDir() string {
	return filepath.Join(getWwwDir(), "software")
}

func getArticlesDir() string {
	return filepath.Join(getWwwDir(), "articles")
}

func getExtremeOptDir() string {
	return filepath.Join(getWwwDir(), "extremeoptimizations")
}

var filesPerDir = make(map[string][]string)

// in logs I saw a url that was correct but had "&foo" and other garbage appended
// to it. Redirect to the best matching file in the directory (if there is a file
// that is a prefix of the file that was asked for)
// returns true if redirected
func redirectIfFoundMatching(w http.ResponseWriter, r *http.Request, dir, fileName string) bool {
	var files []string
	ok := false
	if files, ok = filesPerDir[dir]; !ok {
		files = listFilesInDir(dir, true)
		n := len(dir) + 1
		for i, f := range files {
			files[i] = f[n:]
		}
		//logger.Noticef("files in %q: %v", dir, files)
		filesPerDir[dir] = files
	}
	for _, f := range files {
		if strings.HasPrefix(fileName, f) {
			if fileName == f {
				return false
			}
			diff := len(fileName) - len(f)
			url := r.URL.Path
			url = url[:len(url)-diff]
			logger.Noticef("serveFileFromDir(): redirecting %q => %q", r.URL.Path, url)
			http.Redirect(w, r, url, 302)
			return true
		}
	}
	return false
}

func serveFileFromDir(w http.ResponseWriter, r *http.Request, dir, fileName string) {
	if redirectIfFoundMatching(w, r, dir, fileName) {
		return
	}
	filePath := filepath.Join(dir, fileName)
	if pathExists(filePath) {
		//logger.Noticef("serveFileFromDir(): %q", filePath)
		http.ServeFile(w, r, filePath)
	} else {
		//logger.Noticef("serveFileFromDir() file %q doesn't exist, referer: %q", fileName, r.Referer())
		serve404(w, r)
	}
}

// url: /static/*
func handleStatic(w http.ResponseWriter, r *http.Request) {
	if redirectIfNeeded(w, r) {
		return
	}
	file := r.URL.Path[len("/static/"):]
	serveFileFromDir(w, r, getStaticDir(), file)
}

// url: /css/*
func handleCSS(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/css/"):]
	serveFileFromDir(w, r, getCSSDir(), file)
}

// url: /js/*
func handleJs(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/js/"):]
	serveFileFromDir(w, r, getJsDir(), file)
}

// url: /gfx/*
func handleGfx(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/gfx/"):]
	serveFileFromDir(w, r, getGfxDir(), file)
}

// url: /software*
func handleSoftware(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if url == "/software" || url == "/software/" || url == "/software/index.html" {
		serveFileFromDir(w, r, getSoftwareDir(), "index.html")
		return
	}
	if redirectIfNeeded(w, r) {
		return
	}
	file := r.URL.Path[len("/software/"):]
	serveFileFromDir(w, r, getSoftwareDir(), file)
}

// url: /favicon.ico
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	serveFileFromDir(w, r, getStaticDir(), "favicon.ico")
}

// url: /contactme.html
func handleContactme(w http.ResponseWriter, r *http.Request) {
	serveFileFromDir(w, r, getStaticDir(), "contactme.html")
}

// url: /robots.txt
func handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	serveFileFromDir(w, r, getWwwDir(), "robots.txt")
}

// url: /extremeoptimizations/
func handleExtremeOpt(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/extremeoptimizations/"):]
	if file == "" {
		file = "index.html"
	}
	serveFileFromDir(w, r, getExtremeOptDir(), file)
}

// url: /articles/*
func handleArticles(w http.ResponseWriter, r *http.Request) {
	if redirectIfNeeded(w, r) {
		return
	}
	url := r.URL.Path
	if url == "/articles/" || url == "/articles/index.html" {
		serveFileFromDir(w, r, getStaticDir(), "documents.html")
		return
	}
	file := r.URL.Path[len("/articles/"):]
	serveFileFromDir(w, r, getArticlesDir(), file)
}
