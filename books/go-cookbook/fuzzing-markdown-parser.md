---
Id: 23
Title: Fuzzing Markdown parser written in Go
Date: 2018-01-31T15:32:55-08:00
Tags: go, programming
Format: Markdown
HeaderImage: gfx/headers/header-21.jpg
Collection: go-cookbook
Description: Fuzzing markdown parser written in Go with go-fuzz
---

I finished working on [Markdown parser](https://github.com/gomarkdown/markdown) for Go.

To make sure it's robust (no crashes, no hangs) I decided to fuzz it.

## What is fuzzing?

Parsing text or binary formats is notoriously tricky.

It's easy to make mistakes that can lead to security exploits.

[Fuzzing](https://security.googleblog.com/2016/08/guided-in-process-fuzzing-of-chrome.html) is one technique used to combat this.

The idea is simple:
* magic algorithm generates randomized input data
* we feed this data to parsing code
* code instrumentation detects if parsing code crashes or overwrites memory

## Fuzzing in Go is simple

Writing a fuzzer for C or C++ code is hard.

For Go it's trivial thanks to [go-fuzz](https://github.com/dvyukov/go-fuzz).

It took me about an hour to write a fuzzer for my markdown parser.

## Fuzzing step-by-step

First we need to get the tools:
```bash
$ go get github.com/dvyukov/go-fuzz/go-fuzz
$ go get github.com/dvyukov/go-fuzz/go-fuzz-build
```

My library has `ast := markdown.Parse(markdown []byte, p *parser.Parser)` function which takes markdown text as input and parses it into an abstract syntax tree.

I wrote `fuzz.go`:
```go
// +build gofuzz

package markdown

// Fuzz is to be used by https://github.com/dvyukov/go-fuzz
func Fuzz(data []byte) int {
	Parse(data, nil)
	return 0
}
```

This implements `Fuzz(data []byte)` function, which is an API that go-fuzz expects.

Notice `// +build gofuzz` line.

It means that this code will not be compiled by default but only when compiled with `gofuzz` build tag.

Having the fuzzer generate completely random data is better than nothing but it's much better if fuzzer can be seeded with sample valid content.

As it happens, I already had markdown files used for tests in `testdata` directory.

I used those to initialize the fuzzer:

```bash
# create a working directory for the fuzzer
$ mkdir -p fuzz-workdir

# copy the files that seed fuzzing to corpus sub-directory of working directory
$ mkdir -p fuzz-workdir/corpus
$ cp testdata/*.text fuzz-workdir/corpus

# generate the fuzzing program. This compiles fuzz.go we wrote earlier
# generates fuzzer executable and markdown-fuzz.zip that packages
# data to drive fuzzing process
$ go-fuzz-build github.com/gomarkdown/markdown
```

The last `go-fuzz-build` step can take a while.

Finally we start fuzzing process:
```bash
$ go-fuzz -bin=./markdown-fuzz.zip -workdir=fuzz-workdir
```

This runs for as long as you let it. Fuzzing works by generating random input data so it can go forever.

As long as you don't delete working directory, you can stop and re-start `go-fuzz` and it'll pick up where it left off.

Found crashes will be logged in `fuzz-workdir/crashers` directory.

I'm a fan of automation so I wrote [`fuzz.sh`](https://github.com/gomarkdown/markdown/blob/master/s/fuzz.sh) script.

That way I can re-run it easily after making changes to the parser to verify I didn't introduce bugs.

The fuzzer has advanced options like fuzzing using multiple machines.

## The results

Fuzzing is magic. My library is derived from very popular and widely used blackfriday library (which in turn was derived from a popular C library) and yet the fuzzer found 3 separate issues that either crashed the parser or entered infinite loop.

They [are](https://github.com/gomarkdown/markdown/commit/5d96569c5a0d3cd46d961eddbb61e936e627774c) [now](https://github.com/gomarkdown/markdown/commit/e0fc813169b926a2182bc6554888eb37d12261f7) [fixed](https://github.com/gomarkdown/markdown/commit/5dd4b50fe81eda60f173e242ece05f24c5cc5cec) and covered by unit tests.

Go West and Fuzz!
