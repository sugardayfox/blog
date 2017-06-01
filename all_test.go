package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testShortenID(t *testing.T, n int) {
	s := shortenID(n)
	n2 := unshortenID(s)
	assert.Equal(t, n, n2)
}

func testGzip(t *testing.T, path string) {
	d, err := ioutil.ReadFile(path)
	assert.Nil(t, err)

	dstPath := path + ".gz"
	err = gzipFile(dstPath, path)
	defer os.Remove(dstPath)
	assert.Nil(t, err)
	r, err := openFileMaybeCompressed(dstPath)
	assert.Nil(t, err)
	defer r.Close()
	var dst bytes.Buffer
	_, err = io.Copy(&dst, r)
	assert.Nil(t, err)
	d2 := dst.Bytes()
	assert.Equal(t, d, d2)
	os.Remove(dstPath)
}

func TestGzip(t *testing.T) {
	testGzip(t, "visitor_analytics.go")
}
func TestShortenId(t *testing.T) {
	testShortenID(t, 1404040)
	testShortenID(t, 0)
	testShortenID(t, 1)
	testShortenID(t, 35)
	testShortenID(t, 36)
	testShortenID(t, 37)
	testShortenID(t, 123413343)
}
