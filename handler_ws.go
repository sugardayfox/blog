package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

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

// AddWatch adds watch for a file
func AddWatch(path string) chan struct{} {
	mu.Lock()
	defer mu.Unlock()

	for _, wf := range watchedFiles {
		if path == wf.path {
			return nil
		}
	}

	c := make(chan struct{})
	log.Printf("AddWatch: '%s'\n", path)
	wf := watchedFile{path, c}
	watchedFiles = append(watchedFiles, wf)
	return c
}

// RemoveWatch reoves watch for a file
func RemoveWatch(c chan struct{}) {
	mu.Lock()
	defer mu.Unlock()

	for i, w := range watchedFiles {
		if w.c == c {
			fmt.Printf("removed watching of %s\n", w.path)
			watchedFiles = append(watchedFiles[:i], watchedFiles[i+1:]...)
			return
		}
	}
}

// NotifyFileChanges sends a notification about changed file
func NotifyFileChanges(ev fsnotify.Event) {
	path := ev.Name
	if isTmpFile(path) {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, w := range watchedFiles {
		if strings.HasSuffix(path, w.path) {
			log.Printf("NotifyFileChanges: notify about %s, event: %s\n", path, ev.String())
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
		}
	}
	//watcher.Close()
}

func reloadArticle(article *Article) {
	for i, a := range store.articles {
		if a == article {
			log.Printf("reloading %s\n", article.Path)
			newArticle, err := readArticle(a.Path)
			if err != nil {
				log.Printf("reloading %s failed with %s\n", a.Path, err)
				return
			}
			store.articles[i] = newArticle
			return
		}
	}
}

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func fileWatcher(article *Article, wsChan chan struct{}, chClose chan struct{}) {
	c := AddWatch(article.Path)
	if c == nil {
		return
	}
	defer RemoveWatch(c)
	for {
		select {
		case <-c:
			reloadArticle(article)
			wsChan <- struct{}{}
		case <-chClose:
			return
		}
	}
}

func wsWriter(ws *websocket.Conn, wsChan chan struct{}, chClose chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case <-ticker.C:
			log.Print("sending ws ping")
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Print(err)
				return
			}
		case <-wsChan:
			err := ws.WriteMessage(websocket.TextMessage, []byte{})
			if err != nil {
				log.Printf("wsWrite: ws.WriteMessage failed with '%s'\n", err)
				return
			}
		case <-chClose:
			return
		}
	}
}

// serveWs receives a file name from a websocket client and relays to it
// all the notifications about changes to this file.
func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Printf("serveWs: new connection\n")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Printf("serveWs: upgrader.Upgrade() failed with '%s'\n", err)
		}
		return
	}

	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	wsChan := make(chan struct{})
	chClose := make(chan struct{})
	defer func() {
		chClose <- struct{}{}
	}()
	go wsWriter(ws, wsChan, chClose)
	for {
		msgType, p, err := ws.ReadMessage()
		if err != nil {
			log.Printf("serveWs: ws.ReadMessage() failed with '%s'\n", err)
			return
		}
		log.Printf("Got ws msg type: %d s: '%s'\n", msgType, string(p))
		uri := string(p)
		articleInfo := articleInfoFromUrl(uri)
		if articleInfo == nil {
			log.Printf("serveWs: didn't find article for uri %s\n", uri)
			continue
		}
		article := articleInfo.this
		go fileWatcher(article, wsChan, chClose)
	}
}
