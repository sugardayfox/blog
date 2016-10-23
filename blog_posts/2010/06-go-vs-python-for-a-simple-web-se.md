Id: 204001
Title: Go vs. Python for a simple web server
Tags: python,go
Date: 2010-06-13T22:13:12-07:00
Format: Markdown
--------------
When [Go language](http://golang.org/) was announced, I liked what I
saw. However, to make an informed opinion about a programming language
you really have to implement something in it, and that’s what I did.

The app
=======

I re-wrote a web server applications originally built in Python, using
[Tornado web framework](http://www.tornadoweb.org/).

The app is relatively simple: it implements 5 REST points, each
returning a JSON-formatted data that it gets by executing other
applications and converting their output to a data structure.

Concurrency
===========

Go version used built-in HTTP library (there are other options, like
[web.go](http://github.com/hoisie/web.go)).

If there’s one feature where Go easily beats Python it’s concurrency.
Mainstream languages, including Python, use threads for concurrency and
locking for protecting shared data, which is error prone. In addition,
threading in Python performs poorly due to GIL.

In web serving context, thread-per-connection way of handling HTTP
requests is resource intensive so efficient Python web servers (like
Tornado) are event driven.

The downside of that is that at any given time Tornado can process only
one HTTP request. If that processing takes a long time, it blocks all
other requests.

To fix the single-threaded nature of event driven web servers, people
run several server processes behind an efficient, proxying web server
like nginx, which increases the complexity of the overall solution.

In Go, writing concurrent code is more natural and efficiency of
goroutines is much better than threads. That makes
goroutine-per-connection way of handling HTTP requests more applicable
even in heavily stressed applications.

Source code size
================

Go is remarkably concise for a strongly typed language with C-like
syntax. Python code was about 1000 lines and Go equivalent about 1300
lines. I consider that a great result.

There are some things, like sorting vectors of structs, that require
longer code than in Python (one case where dynamic typing helps) but
overall you can write very compact code in Go.

Speed and memory efficiency
===========================

I haven’t measured. Generally speaking, Go is compiled to assembly code
and its data structures don’t suffer Python’s terrible per-object
overhead, so it should be more efficient in both speed and memory usage.

Currently Go’s speed is handicapped by a weak garbage collector (an
issue known to Go’s developers who claim this is actively being worked
on)

Other things I like about Go
============================

Go’s syntax is very C-like but it cleverly gets rid of one source of
noise in C: semicolons. They are still there but they can be omitted in
most cases (Go compiler figures out where a semicolon is supposed to
be). It makes Go code look cleaner.

I like how classes were replaced by interfaces. They especially shine in
i/o code and layered streams (e.g. a compressed stream that wraps
encrypted stream that wraps stream from a file or network or memory).
They’re easy to use from the client point of view and it’s easy to
implement a new type of wrapper (e.g. using a new encryption algorithm
or new compression scheme).

Another nice property of Go is how it minimizes indentation (compared to
e.g. C++ code or Java code).

C and C++ suffer from big variations in programming style which makes
combining two code bases written in different styles or reading code in
a different style that you’re used to hard.

Python solved this problem well by requiring white space for indentation
instead of free-form braces.

Go fixes this problem by convention: while one can be as inconsistent as
in C, go includes `gofmt` program, which will reformat the code in a way
consistent with official Go style. Go’s standard library code uses this
style and most developers writing their own libraries seem to have
adopted that as well, to play nice with the ecosystem.

Go’s downsides
==============

Go’s downside is lack of maturity and weak Windows support. There are
still rather serious bugs being fixed in the compiler and the runtime,
built-in libraries aren’t as complete as in older implementations like
Perl, Python or Ruby, there’s much less third party libraries.

In defense of Go, you can hardly expect miracles. Go is a very young
implementation and seems to be more mature that other language
implementations were at that age.

Conclusion
==========

Overall, the experience was positive and I would seriously consider
using Go for writing server side of web applications. It provides a
right balance between speed and memory efficiency of the running code
(much closer to C than to e.g. Python) and speed of development (much
closer to Python than to C).

It’s also reasonable to expect that it’ll only get better in the future.
While Google doesn’t devote a huge team to Go, those who are on the team
have demonstrated in the past that they know how to get things done.

Go is used internally at Google, so it won’t be abandoned.

As far as open-source projects go, this one is well managed i.e. Go team
encourages outside code contributions and is quite speedy in reviewing
and committing external patches.
