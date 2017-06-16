Id: 985008
Title: SumatraPDF 2.0 released
Tags: sumatra,releasenotes
Date: 2012-04-02T14:54:09-07:00
Format: Markdown
--------------
[SumatraPDF](https://www.sumatrapdfreader.org/free-pdf-reader.html) is a
multi-format PDF, XPS, ebook (MOBI), DjVu etc. reader for Windows and
[we](http://www.ohloh.net/p/4623/contributors) are pleased to announce
version 2.0.

The biggest change in this version is that we can now read ebooks in
mobi format. Mobi is a format developed for MobiPocket Reader and made
popular by Amazon’s Kindle. Unfortunately, eBooks purchased from Amazon
are protected by DRM so SumatraPDF (or any other third-party reader)
cannot read them. The good news is that there are other sources of
ebooks formatted in mobi format.

The UI for reading ebooks is different than the UI used for other
documents. It’s not hard to notice that the inspiration from the UI came
from Kindle’s PC application.

There were other changes as well:

We can now open CHM documents from network drives.

After selecting an area with the mouse, the area can be copied to the
clipboard as an image with a right-click context menu.

Sumatra has always been a small and fast program and because we applied
[extreme size optimization
techniques](http://code.google.com/p/sumatrapdf/source/browse/trunk/src/ucrt/readme.txt),
it’s smaller than ever. If we didn’t apply our extreme size-reduction
techniques, the installer would be bigger by 9% (which is around 400 kB)
(as a bonus, the size-reducing code we developed is available to other
programmers under liberal BSD license).

And as always there are many smaller improvements: even better PDF
support etc.
