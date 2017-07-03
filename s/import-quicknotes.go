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

func getURLPrefix(uri string) string {
	urlPrefix := "https://quicknotes.io/raw/n/"
	if strings.HasPrefix(uri, urlPrefix) {
		return urlPrefix
	}
	urlPrefix = "http://quicknotes.io/raw/n/"
	if strings.HasPrefix(uri, urlPrefix) {
		return urlPrefix
	}
	msg := fmt.Sprintf("url '%s' should start with '%s' or '%s'", uri, "https://quicknotes.io/raw/n/", "http://quicknotes.io/raw/n/")
	panic(msg)
}

// given a date in format "YYYY-MM-DD", return it format "2017-06-21T18:47:07Z"
// empty string returns empty string
// panic if non-empty but doesn't fit the expected format
func validateDate(s string) string {
	s = strings.TrimSpace(s)
	t, err := time.Parse("2006-01-02", s)
	u.PanicIfErr(err)
	return t.Format(time.RFC3339)
}

func main() {
	path := filepath.Join("articles", "from-quicknotes.txt")
	lines, err := u.ReadLinesFromFile(path)
	u.PanicIfErr(err)
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

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// a line could be in format: ${url} ${published_on}
		parts := strings.Split(line, " ")
		u.PanicIf(len(parts) > 2)
		uri := parts[0]
		urlPrefix := getURLPrefix(uri)
		d, err := httpGet(uri)
		u.PanicIfErr(err)
		if len(parts) == 2 {
			publishedOnStr := validateDate(parts[1])
			d2 := []byte("PublishedOn: " + publishedOnStr + "\n")
			d = append(d2, d...)
		}

		name := strings.TrimPrefix(uri, urlPrefix)

		parts = strings.SplitN(name, "-", 2)
		noteID := parts[0]
		fileName := noteIDToFileName[noteID]
		//fmt.Printf("noteID: %s, fileName: %s\n", noteID, fileName)
		if fileName == "" {
			fileName = name + ".md"
		}
		path := filepath.Join(dir, fileName)

		err = ioutil.WriteFile(path, d, 0644)
		u.PanicIfErr(err)
		fmt.Printf("Wrote '%s' as '%s', %d bytes\n", uri, path, len(d))
	}
}
