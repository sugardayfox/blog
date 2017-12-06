Id: 486053
Title: Experience porting 4k lines of C code to go
Tags: go,programming
Date: 2011-06-01T23:40:38-07:00
Format: Markdown
--------------
I’ve been recently interested in [Go](https://golang.org), a new
programming language.

The best way to learn a language is to use it in a small but real
project.

I needed a program that generates HTML from textile format. The way of
least resistance would be to implement it in Python, as I’ve coded a
similar thing in the past. Instead I decided to use this as an
opportunity to get more familiar with Go.

The biggest part of the problem was the textile to HTML conversion.
There is no existing Go code for that so I decided to port
[upskirt](https://github.com/tanoku/upskirt) C library, as it does the
job in the most performant way (it has a hand-written, disciplined
parser as opposed to most other solutions that just throw cryptic
regular expression at the problem).

The bottom line is: porting C code, at least in this case, was fast,
boring, mechanical process (and that is a good thing).

Go’s syntax is heavily inspired by C. The differences that I’ve
encountered most frequently:

-   syntax for declaring variables is different (and better)
-   while keyword is missing, replaced by a more versatile for syntax
-   function declaration syntax is different

Fortunately, the transformations were simple and mostly mechanical.

It took me just few days to manually translate around 4000 lines of C
code into \~3200 lines of Go code.

Saving \~800 lines of code (20%) is good, but the interesting part is:
where do the savings come from?

The core parsing/html generating logic didn’t shrink much. The big
savings came from the fact that Go has a built-in growable array type
and upskirt C code had to spend 924 lines re-implementing that in C.

An unexpected advantage of Go was its safety. The C code implements
parsing by partying on char \* pointers. Such code is notorious for
causing lots of subtle, hard to test for bugs. Go doesn’t allow this
kind of pointer arithmetic and instead provides slices, which are a view
into an underlaying array.

Slices provide out-of-bounds checks. Just by recoding in Go I found
out-of-bounds access bug in the [original C
code](https://github.com/tanoku/upskirt/issues/24).

Thanks to similarity of Go and C syntax, porting algorithmic code from C
is simple.

All things considered, Go is quickly becoming my new preferred language
(taking the crown away from Python). It combines the good attributes of
Python (lightweight syntax, garbage collection) with good attributes of
C (fast execution thanks to compilation to native code and programmer’s
control over memory layout) and adds some unique capabilities of its own
(concurrency via gorutines and channels).

BTW: if you want markdown implementation for Go, use
[blackfriday](https://github.com/russross/blackfriday). It’s also direct
port of upskirt and I abandoned my port in favor of contributing to
blackfriday, since it was slightly ahead and there’s no need for two
nearly identical projects.
