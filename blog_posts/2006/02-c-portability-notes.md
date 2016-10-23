Id: 1293
Title: C portability notes
Tags: c,programming
Date: 2006-02-11T16:00:00-08:00
Format: Markdown
--------------
Using `#ifdef` is a sad necessity for portable C\\C++ programming. This is a list
of a few symbols that are often used for writing really portable programs (those
that run on Unix\\Windows\\MacOS\\Palm OS etc.). This list is meant as a quick
reference, not an in-depth explanation.

`_WIN32_WCE` - defines Windows CE version (i.e. if defined, this is Windows CE)

`_WIN32` and `_WINDOWS` - usually mean compilation for Windows

`__CYGWIN__` - compiled by gcc under CYGWIN (Unix emulation layer for Windows)

`__MINGW32__` - compiled by gcc under Mingw32 - Unix portability layer for Windows

`_UNICODE` - on Windows means that TCHAR is WCHAR (or wchar_t) i.e. 16-bit unicode character. On Windows CE this is the only option.

`__BORLANDC__` - set if compiled with Borland C compiler

`_MSC_VER` - defines a version of Microsoft C compiler

`__GNUC__` - defined when using gcc

`NDEBUG` - defined if this is a release (i.e. not DEBUG) build

`DEBUG` - often defined for debug builds (opposite of NDEBUG)


`_PALM_OS` - compilation for Palm OS
