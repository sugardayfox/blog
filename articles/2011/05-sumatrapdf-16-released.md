Id: 498008
Title: SumatraPDF 1.6 released
Tags: sumatra,releasenotes
Date: 2011-05-30T21:42:25-07:00
Format: Markdown
--------------
[SumatraPDF developers](http://www.ohloh.net/p/4623/contributors) are
quite pleased to announce 1.6 release of
[SumatraPDF](https://www.sumatrapdfreader.org/free-pdf-reader.html), a PDF,
XPS, DjVu, CBZ and CBR reader for Windows.

In this release we’ve added support for [DjVu](http://djvu.org/) file
format.

When no document is open, we display a list of frequently read document
as thumbnails. This functionality is inspired by new tab page in Chrome.

We’ve added support for displaying Postscript documents. This requires
recent Ghostscript version to be already installed - we don’t bundle it
ourselves.

We’ve added support for displaying a folder containing images: drag the
folder to SumatraPDF window

We now support clickable links and a Table of Content for XPS documents.

We’ve added printing progress and allow canceling printing process.

We’ve added Print toolbar button.

Experimental: we’ve added previewing of PDF documents in Windows Vista
and 7. This creates thumbnails and displays documents in Explorer’s
Preview pane. Needs to be explicitly selected during install process.
We’ve had reports that it doesn’t work on 64-bit Windows which is why we
call it experimental.

This is how “frequently read” list looks like:

![](=//kjkpub.s3.amazonaws.com/blog/sumatra/sum-shot-03-small.png)
