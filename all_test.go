package main

import (
	_ "fmt"
	"testing"
)

func testShortenId(t *testing.T, n int) {
	s := ShortenID(n)
	n2 := UnshortenID(s)
	if n != n2 {
		t.Fatalf("'%d' != '%d', shortened = %q", n, n2, s)
	}
}

func TestShortenId(t *testing.T) {
	testShortenId(t, 1404040)
	testShortenId(t, 0)
	testShortenId(t, 1)
	testShortenId(t, 35)
	testShortenId(t, 36)
	testShortenId(t, 37)
	testShortenId(t, 123413343)
}
