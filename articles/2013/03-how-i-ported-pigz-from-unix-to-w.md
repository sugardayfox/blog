Id: 4
Title: How I ported pigz from Unix to Windows
Tags: programming
Date: 2013-03-20T11:25:36-07:00
Format: Markdown
--------------
I just finished [porting](https://github.com/kjk/pigz) [pigz](http://zlib.net/pigz/) from Unix to Windows (you can download [pre-compiled binaries](/software/pigz-for-windows.html)). This article describes how I did it.

Pigz was clearly written with Unix in mind, with no thought given to cross-platform portability.

Thankfully, it's a relatively simple, command-line program that sticks to using standard C library.

## Porting pthreads

Pigz uses pthreads for threading. Porting pthreads code to Windows would be a nightmare. Lucky me: someone already did all the hard work and [implemented pthreads APIs](https://github.com/GerHobbelt/pthread-win32) on top of Windows API, in only 20.000 lines of code. It seems to Just Work.

## Porting dirent

Another Unix-only API that pigz uses is [dirent.h](http://pubs.opengroup.org/onlinepubs/007908799/xsh/dirent.h.html), for reading the content of directories. I was lucky again: someone created a [Windows port of dirent APIs](http://www.two-sdg.demon.co.uk/curbralan/code/dirent/dirent.html).

## Misc fixes

There are few functions that pigz uses that are present in Visual Studio's C library, but under a different name (e.g. `stat` is `_stat`, `fstat` is `_fstat` etc.). This is easily fixed with a `#define`.

Visual C doesn't define `ssize_t` and `PATH_MAX`. A `typedef` here, a `#define` there solves those problems.

I consolidated all such fixes in a single [wincompat.h](https://github.com/kjk/pigz/blob/master/win32/wincompat.h) file.

Pigz `#include`s some .h files that are only available under Unix. I used `#ifndef _WIN32` around those.

## Build system

Pigz uses standard Unix build tools: gcc and make.

While there are ports of GNU make to Windows and gcc-based compilers that can generate windows binaries (mingw), most Windows developers prefer using Visual Studio IDE.

Creating Visual Studio project files from scratch is time consuming and annoying, especially when you want to support multiple versions of Visual Studio.

[Premake](http://industriousone.com/premake) makes it easier. It's a meta-build system: you write a [text file](https://github.com/kjk/pigz/blob/master/premake4.lua) that defines the project, which files to compile, compilation flags etc. and premake generates Visual Studio files from that description.

Premake supports Visual Studio 2008 and 2010 (and 2012 supports 2010 project files via conversion).

## Conclusion

As far as porting goes, this one was easy thanks to existing efforts that created compatibility shims for the hard parts.

Premake is an interesting tool that allows to save time creating and maintaining Visual Studio projects.
