package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type watchedFile struct {
	path string
	c    chan struct{}
}

var (
	mu           sync.Mutex
	watchedFiles []watchedFile
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func AddWatch(path string) chan struct{} {
	c := make(chan struct{})
	mu.Lock()
	wf := watchedFile{path, c}
	watchedFiles = append(watchedFiles, wf)
	mu.Unlock()
	return c
}

func RemoveWatch(c chan struct{}) {
	mu.Lock()
	for i, w := range watchedFiles {
		if w.c == c {
			fmt.Printf("removed watching of %s\n", w.path)
			watchedFiles = append(watchedFiles[:i], watchedFiles[i+1:]...)
			break
		}
	}
	mu.Unlock()
}

func NotifyFileChanges(ev fsnotify.Event) {
	path := ev.Name
	if isTmpFile(path) {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, w := range watchedFiles {
		if strings.HasSuffix(path, w.path) {
			fmt.Printf("NotifyFileChanges: notify about %s, event: %s\n", path, ev.String())
			select {
			case w.c <- struct{}{}:
			default:
			}
		}
	}
}

func watchChanges(watcher *fsnotify.Watcher) {
	for {
		select {
		case ev := <-watcher.Events:
			NotifyFileChanges(ev)
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func startWatching() {
	if inProduction {
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("fsnotify.NewWather() failed with %s\n", err)
		return
	}

	go watchChanges(watcher)

	dirs := store.GetDirsToWatch()
	dirs = append(dirs, "blog_posts")
	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			fmt.Printf("watcher.Add() for %s failed with %s\n", dir, err)
			return
		} else {
			//fmt.Printf("added watching for dir %s\n", dir)
		}
	}
	//watcher.Close()
}

func reloadArticle(article *Article) {
	for i, a := range store.articles {
		if a == article {
			fmt.Printf("reloading %s\n", a.Path)
			newArticle, err := readArticle(a.Path)
			if err != nil {
				fmt.Printf("reloading %s failed with %s\n", a.Path, err)
				return
			}
			store.articles[i] = newArticle
			return
		}
	}
}

// serveWs receives a file name from a websocket client and relays to it
// all the notifications about changes to this file.
func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Print(err)
		return
	}
	uri := string(p)
	articleInfo := articleInfoFromUrl(uri)
	if articleInfo == nil {
		fmt.Printf("serveWs: didn't find article for uri %s\n", uri)
		return
	}

	article := articleInfo.this
	fmt.Printf("serveWs: started watching '%s' for url '%s'\n", article.Path, uri)

	c := AddWatch(article.Path)
	defer RemoveWatch(c)
	done := make(chan struct{})

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("serveWs: closing for %s\n", uri)
				close(done)
				break
			}
		}
	}()

loop:
	for {
		select {
		case <-c:
			reloadArticle(article)
			err := conn.WriteMessage(websocket.TextMessage, nil)
			if err != nil {
				log.Print(err)
				break loop
			}
		case <-done:
			break loop
		}
	}
}
