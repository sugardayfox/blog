---
Id: wN9R
Title: Experience porting 4.5k loc of C to Go (Facebook's CSS flexbox implementation Yoga)
Format: Markdown
Tags: for-blog, draft
CreatedAt: 2017-07-23T05:35:08Z
UpdatedAt: 2017-08-02T18:02:19Z
PublishedOn: 2017-08-02
HeaderImage: gfx/headers/header-18.jpg
Collection: go-cookbook
Description: Porting Yoga (Facebook's CSS flexbox implementation) from C to Go.
---
I'm working on desktop GUI apps for Windows in Go.

As it uses native win32 controls I don't have to implement my own buttons, list boxes etc.

However, win32 doesn't have anything to help you layout things on the screen. Manual positioning is painful.

I started looking into a writing a layout engine. There are many systems to be inspired by.

I did some work in Winforms so I could copy their ideas.

I could look into copying Android layout logic.

I could look into implementing constraint-based layout like in iOS.

The option I liked the best is CSS flexbox. It's not the simplest, it's not the most complicated, but the biggest win is that many people, including myself, are already familiar with it.

I started looking for an implementation in Go. I found one in [shiny](https://github.com/golang/exp/tree/master/shiny/widget/flex) project, but it's simplified implementation.

The most complete implementation is Facebook's [Yoga](https://github.com/facebook/yoga), but it's written in C (with bindings to many languages).

I could write cgo bindings but I prefer my Go pure so I decided to manually port C code to Go, all 4.5k lines of it.

The result is [flex](https://github.com/kjk/flex) package.

Here's how porting happened.

## Phase zero - picking a name for the package

It was tempting to pick a name that would be some combination of `yoga` and `go` e.g. `goyoga` or `yoga-go`.

That would be a bad idea. Including `go` in package or repository name is bad (but unfortunately happens too often).

My rules for a good name for a Go library are:

* it's short
* it's descriptive
* package name matches repository name

Therefore [github.com/kjk/flex](github.com/kjk/flex) was born.

I could use `yoga`, to better higlight connection with original project, but it's not descriptive.

I could use `flexbox` to better higlight connection with [CSS flexbox](https://www.w3.org/TR/css-flexbox-1/), but it's a bit long.

## Phase one - getting Go port to compile

You can't port 4.5k lines of code in one sitting so it was important to have a strategy that allows for making incremental progress.

The important part was keeping track of porting progress. For that I cloned Yoga's repository. After converting a piece of C code to Go I would delete it from my fork. That way I would keep track of code that still needs to be ported.

It's important to get an early boost so I started with [the simplest parts](https://github.com/kjk/flex/commit/4dabacad1ba37403e8b7380b2ccb6c2ed6c8586a): structures, enums.

Porting process was boring but relatively uneventful.

Some porting comments:
* Go and C have reversed order of declarations. Reversing it manually is boring and error prone. There were many repeated declarations that I could do with a simple search-and-replace, e.g `YGNodeRef node` => `node *YGNode`
* search-and-replace `->` to `.` as Go uses `.` for both cases (something that C++-28 should adopt)
* Go doesn't support ternary operator and Yoga's developers are infatuated with it. That was one part where I had to be extra careful as it was more than mechanical change
* Go doesn't support `xor` (^) and `or` (|) operator on bool. As it turns out, they are not neccessary. `x ^ y` is the same as `x != y`. `a = a | b` can be written as: `if !a && b { a = true }`
* `switch` statement needs attention as in C `case` falls-through by default
* Yoga uses `float` for coordinates and uses a few functions like `fmaxf` etc. Go only implements `float64` versions of math functions. Missing functions were easy to write.
* Yoga code was relatively easy to port partly because of lack of string handling, which is where Go and C differ a lot

It took me 2 days to port C code to Go, get it compiling and passing some tests. But some tests were failing.

## Phase two - fixing bugs

Yoga has a lot of tests which gave me confidence that when the port is done,  I'll be able to verify its correctness by porting the tests.

I ported the code, I started porting the tests and they were failing.

I was in a pickle.

2 days was enough to mechanically port the code but not even close enough to understand it.

As far as I can tell you can't step through tests in the debugger, so I converted a failing test into a test program.

I was pleasently surprised that debugger in Visual Studio Code works decently.

I found one porting mistake by stepping through the code but fixing it didn't fix test failures.

I spent a bit more time stepping through the code but not knowing what results to expect I decided that it's not a promising approach.

I decided to re-check every line of code.

Initially I ported the code in random order so I used re-checking time to also re-arrange the code in the same order as C code. That will help re-porting Yoga improvments in the future.

It was predictably boring. It took another day, I found 2 more porting bugs.

The tests were still failing and I was a bit stuck.

I did another pass in the debugger and fortunately an inspiration struck.

Yoga uses `NaN` value as `undefined`. I noticed that `fmaxf` function that I wrote returned `NaN` when any argument was `NaN`.

I compared my `fmaxf` implementation with C implmentation and it turns out that if one of the arguments is `NaN` but the other one isn't, it should return non-NaN value.

It was easy to fix and that fixed all the tests.

## Phase three - Go-ifying the API

The code was working but it was far from idiomatic Go code.

Next phase was tweaking the API to be more Go-like.

The majority of changes were:

* making struct fields public
* removing `YG` prefix as packages obviate the need for that
* renaming accessor functions to methods (from `func NodeFoo(node *Node)` => `func (node *Node) Foo()`

It's mechanical, boring process.

Some transformations could be done with regex search & replace. The rest with manual labor.

Having lots of tests really helped in making sure that the code is still correct.

## The status

The port is finished. It works and passes all Yoga tests.

The API is still a bit awkward by Go standards. Too many methods, not enough direct access to variables.

Partly it's because I wanted to stay close to C code so that future changes to Yoga can be ported easily.

Partly it's because of implementation choices and the nature of the problem.

There are many flexbox properties, only some are given for any given element. The rest is given default values, which might be undefined.

How to represent that in Go? Ideally, the zero value of the type would represent default value so that a style defintion can be constructed as `&Style{}` and then we can set the properties that are non-default.

Unfortunately, that doesn't work well for properties that are numbers because zero might be a valid value.

Yoga uses `float32` for numbers and `NaN` for undefined value of a CSS property.

Another option is used in `sql` package for types like `sql.NullString` that need to indicate null-abilibility. They are represented as a struct that combines a value and bool field `Valid`. Zero value of bool is false, so by default values start as invalid (undefined).

This isn't a great API either since it makes setting and reading values more awkward.

At the end of the day it's better to have solid, working code with awkward API than not having the code at all. Full flexbox implementation is not trivial to write.

Unless someone writes a pure Go implementation designed from scratch with Go idioms in mind, [flex](https://github.com/kjk/flex) is the best option.

## Notes on automatic translation

Go is so close to C. Wouldn't it be great if there was a program that could take C code and turn it into Go code?

There are few attempts to do that:
* https://github.com/rsc/c2go
* https://github.com/elliotchance/c2go
* https://github.com/cznic/ccgo

[elliotchance/c2go](https://github.com/elliotchance/c2go) seems to be the most promising.

I didn't try any of them as neither seems to be usable yet.

It's an approach worth keeping an eye on for the future but not yet viable today.
