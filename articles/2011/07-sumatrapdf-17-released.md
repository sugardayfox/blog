Id: 575001
Title: SumatraPDF 1.7 released
Tags: sumatra,releasenotes
Date: 2011-07-17T18:28:26-07:00
Format: Markdown
--------------
[SumatraPDF developers](http://www.ohloh.net/p/4623/contributors) are
pleased to announce 1.7 release of
[SumatraPDF](http://blog.kowalczyk.info/software/sumatrapdf), a PDF,
DjVu, XPS, CBZ and CBR reader for Windows.

In this release we’ve added user-defined favorites (i.e. bookmarks). You
can create one or more favorites for a given file, navigate to a
favorite and delete them.

Favorites are accessed either via a menu items in Favorites top-level
menu or displayed as a tree in the sidebar.

We’ve improved support for right-to-left languages, like Arabic.

Logical page numbers are displayed and used, if document provides them
(such as i, ii, iii, etc.).

We allow to restrict SumatraPDF’s features with more granularity; see
[this
document](http://code.google.com/p/sumatrapdf/source/browse/trunk/docs/sumatrapdfrestrict.ini)
for more information.

Command-line argument -named-dest now also matches strings in table of
contents.

We’ve improved support for EPS files (requires Ghostscript)

Installer is now more robust. Previously an installation could fail if a
web browser using Sumatra’s web browser dll was running. Now installer
detects this and will ask to close the browser before proceeding.

Until next release.
