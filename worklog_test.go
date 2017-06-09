package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveHashtags(t *testing.T) {
	tests := []string{
		"  https://blog.figma.com/webassembly-cut-figmas-load-time-by-3x-76f3f2395164. Summary: webassembly is much faster than asm.js. #webassembly  ",
		"https://blog.figma.com/webassembly-cut-figmas-load-time-by-3x-76f3f2395164. Summary: webassembly is much faster than asm.js.",
	}
	n := len(tests) / 2
	for i := 0; i < n; i++ {
		s := tests[i*2]
		exp := tests[i*2+1]
		got := removeHashtags(s)
		assert.Equal(t, exp, got)
	}
}
