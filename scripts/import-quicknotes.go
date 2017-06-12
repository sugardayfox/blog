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
	path := filepath.Join("articles", "from-simplenote.txt")
	lines, err := u.ReadLinesFromFile(path)
	u.PanicIfErr(err)
	for _, url := range lines {
		url = strings.TrimSpace(url)
		d, err := httpGet(url)
		u.PanicIfErr(err)
		fmt.Printf("url: %s, %d bytes\n", lines, len(d))
	}
}
