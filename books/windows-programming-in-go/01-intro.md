---
Id: wpig-01
Title: Introduction to Windows Programming in Go
Format: Markdown
Tags: go
CreatedAt: 2017-07-12
PublishedOn: 2017-09-12
Collection: go-windows
Description: Introduction to Windows Programming in Go.
Status: invisible
---

This books is meant to teach you how to write native Windows programs in Go.

It's mostly about, but not limited to, GUI programs. I'll cover the traditional win32 APIs, not the latest UWP programs.

If you know Windows programming well, you'll learn how to transfer that knowledge to Go.

If you don't know Windows programming at all, I'll show you engouh to start writing real programs.

## Why write Windows programs in Go?

Why write programs in Go as opposed to C#/Winforms, C#/WPF, C++ using directly win32 APIs or some GUI toolkit like Qt?

If you're using those technologies and are happy with them, don't be offended and keep using them.

That being said, Go offers some good trade-offs.

If you know Go and not those other languages, you can leverage your skills in one more domain.

Unlike .NET (Winforms/WPF) Go doesn't require a large framework installed. It builds static executables.

Size of the executables is not as small as C++ but Go is much more productive. And if you bundle a library like Qt, you'll end up bigger in C++ too.

Go's perfromance is somewhere between C++ and C#. It's plenty fast for most use cases.

Go's concurrency support is unparalleled. Thanks to goroutines making GUI programs that don't block is easy.

Finally, you can leverage a considerable richness of Go ecosystem. There are libraries for pretty much everything you need and they are just `go get` away.
