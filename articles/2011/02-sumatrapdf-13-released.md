Id: 406002
Title: SumatraPDF 1.3 released
Tags: sumatra,releasenotes
Date: 2011-02-07T13:26:41-08:00
Format: Markdown
--------------
![Sumatra
screenshot](http://kjkpub.s3.amazonaws.com/blog/sumatra/sum-shot-02-small.png "Sumatra screenshot")
SumatraPDF team is happy to announce a 1.3 release of our free, small,
fast, open source <span class="caps">PDF</span> reader for Windows.

Sumatra was capable of text selection for some time but this feature was
hard to discover (you had to press Ctrl and select the area with a
mouse). We improved text selection and copying by emulating the way a
browser (and Adobe Reader) work: just select text with a mouse and use
Ctrl-C to copy it to a clipboard.

Shift + Left Mouse now scrolls the document and Ctrl + Left mouse
continues to create a rectangular selection (for copying images).

We added more keyboard and mouse shortcuts:

-   'c’ shortcut toggles continuous mode
-   '+’ / '\*’ on the numeric keyboard now do zoom and rotation
-   back/forward mouse buttons for back/forward navigation

We added toolbar icons for Fit Page and Fit Width and updated the look
of toolbar icons

In version 1.2 we introduced a new full screen mode and made it the
default full screen mode. Old mode was still available but not easily
discoverable. To make it more discoverable we’ve added View/Presentation
menu item for new full screen mode and View/Fullscreen menu item for the
old full screen mode.

We [rewrote the
installer](/article/8nqb/Writing-a-custom-installer-for-windows-software.html)

We improved zoom performance and fixed crashiness caused by high zoom
levels.

We improved searching for text to use less memory.

We improved printing.

We’ve updated translations contributed by our
[translators](https://github.com/sumatrapdfreader/sumatrapdf/blob/master/TRANSLATORS)

We updated to latest [mupdf](http://mupdf.com/) code for various
improvements and bugfixes.

We now use [libjpeg-turbo](http://libjpeg-turbo.virtualgl.org/) library
instead of libjpeg, for faster decoding of some PDFs.

We updated [openjpeg](http://www.openjpeg.org/) library to version 1.4
and [freetype](http://freetype.sourceforge.net/) to version 2.4.4.

We fixed 2 integer overflows reported by Stefan Cornelius from Secunia
Research.
