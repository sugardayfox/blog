---
Id: wpig-03
Title: Hello World, or your first Windows Go program
Format: Markdown
Tags: go
CreatedAt: 2017-07-12
PublishedOn: 2017-09-12
Collection: go-windows
Description: Writing "Hello World" program.
Status: invisible
---

This is not a book to teach you about win32 programming. There are existing resources for that, for example:
* http://www.winprog.org/tutorial/
* http://www.functionx.com/win32/Lesson01.htm
* http://www.relisoft.com/Win32/index.htm
* http://win32-framework.sourceforge.net/tutorial.htm
* http://zetcode.com/gui/winapi/

What I aim is to show how to translate that knowledge into Go and how to use the `github.com\lxn\walk` library.

Let's start by transcribing a C++ [Hello World](https://msdn.microsoft.com/en-us/library/bb384843.aspx) application:

**TODO:** screenshot

Very briefly:
* Windows GUI consists of nested windows, represented by opaque type HWND
* each HWND belongs to a unique class name and that class name determines window procedure to be used
* window procedure is called by Windows when our window procedure needs to process a given window message like keyboard event or mouse move. A lot of Windows programming is writing the right code to handle various messages
* a HWND without a parent is called top-level HWND. This is the window
* most other HWND windows have a parent and are child windows
* Windows OS provides a custom set up child windows called common controls. We know their unique class names, Windows implements all the necessary logic (i.e. a window procedure)

