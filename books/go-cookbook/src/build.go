package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	//"template"
)

func parse(r io.Reader) {
	r2 := bufio.NewReaderSize(r, 2048)
	line := 0
	for {
		l, isPrefix, err := r2.ReadLine()
		if l == nil {
			fmt.Printf("We had %d lines\n", line)
			break
		}
		if err != nil {
			fmt.Printf("We had %d lines\n", line)
			fmt.Printf("%+v", err)
			break
		}
		line++
		if isPrefix {
			fmt.Printf("isPrefix on line %d\n", line)
		} else {
			fmt.Printf("%d: %s\n", line, string(l))
		}
		if line > 500 {
			break
		}
	}
}

func main() {
	dstdir := filepath.Join("..", "..", "www", "go-cookbook")
	_, err := os.Stat(dstdir)
	if err != nil {
		fmt.Printf("Dest dir doesn't exist, creating\n")
		os.MkdirAll(dstdir, 0666)
	}
	srcname := "book.txt"
	f, err := os.Open(srcname)
	if f == nil {
		fmt.Printf("Can't open %q\n", srcname)
		os.Exit(1)
	}
	defer f.Close()
	parse(f)
}
