package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/kjk/u"
)

var (
	// Chrome 59
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36"
)

// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
func makeClient() *http.Client {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func makeGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

func httpGet(url string) ([]byte, error) {
	httpClient := makeClient()
	req, err := makeGetRequest(url)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	path := filepath.Join("articles", "from-quicknotes.txt")
	lines, err := u.ReadLinesFromFile(path)
	u.PanicIfErr(err)
	urlPrefix := "https://quicknotes.io/raw/n/"
	dir := filepath.Join("articles", "from-quicknotes")

	// titles might change but we don't want to change file names
	files, err := ioutil.ReadDir(dir)
	u.PanicIfErr(err)
	noteIDToFileName := make(map[string]string)
	for _, fi := range files {
		u.PanicIf(filepath.Ext(fi.Name()) != ".md")
		parts := strings.SplitN(fi.Name(), "-", 2)
		noteID := strings.TrimSuffix(parts[0], ".md")
		noteIDToFileName[noteID] = fi.Name()
		//fmt.Printf("%s => %s\n", noteID, fi.Name())
	}

	for _, url := range lines {
		url = strings.TrimSpace(url)
		if len(url) == 0 {
			continue
		}
		u.PanicIf(!strings.HasPrefix(url, urlPrefix), "url '%s' should start with '%s'", url, urlPrefix)
		d, err := httpGet(url)
		u.PanicIfErr(err)
		name := strings.TrimPrefix(url, urlPrefix)

		parts := strings.SplitN(name, "-", 2)
		noteID := parts[0]
		fileName := noteIDToFileName[noteID]
		//fmt.Printf("noteID: %s, fileName: %s\n", noteID, fileName)
		if fileName == "" {
			fileName = name + ".md"
		}
		path := filepath.Join(dir, fileName)

		err = ioutil.WriteFile(path, d, 0644)
		u.PanicIfErr(err)
		fmt.Printf("Wrote '%s' as '%s', %d bytes\n", url, path, len(d))
	}
}
