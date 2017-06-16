Id: 427002
Title: SumatraPDF 1.4 released
Tags: sumatra,releasenotes
Date: 2011-03-12T16:18:46-08:00
Format: Markdown
--------------
Good news everyone: a new version of
[SumatraPDF](https://www.sumatrapdfreader.org/free-pdf-reader.html),
a free, small, fast, open-source PDF viewer for Windows, is ready.

What’s new in this version?

We’ve added browser plugin for Firefox/Chrome/Opera (Internet Explorer
is not supported).

We’ve added IFilter that enables full-text search of PDF files in
Windows Desktop Search (i.e. search from Windows Vista/7’s Start Menu).

Browser plugin and IFilter are not installed by default so you need to
use options in the installer and check the appropriate checkboxes. You
can uninstall them by re-running the installer and de-selecting the
options.

In 1.3 we’ve improved text selection to use left mouse button. The
downside of that was that left mouse button couldn’t be used for
scrolling when mouse cursor is over text. To compensate for that, in 1.4
you can scroll with right mouse button.

In 1.3 we’ve introduced a new installer. Some of you missed the ability
of choosing a custom installation directory that was not implemented in
1.3. We’ve re-introduced support for non-standard installation directory
in 1.4.

SumatraPDF focuses on reading PDFs and doesn’t implement some more
advanced functionality like filling out forms or making annotations. If
you need this functionality, we’ve made it easy to re-open current
document via File menu in Adobe Reader (if it’s installed). In 1.4 we
also support Foxit and PDF-XChange.

To make SumatraPDF files smaller, we used to compress the executable
with mpress. Unfortunately that caused some anti-virus programs to
falsely report Sumatra as a virus. We no longer compress the executables
that ship with the installer version, so don’t be surprised that the
files are now bigger. The portable, .zip version still ships as a
single, compressed executable (so don’t be surprised if it’s flagged by
some anti-virus software).

We’ve removed -title cmd-line option.

We’ve added support for AES-256 encrypted PDF files, fixed an integer
overflow reported by Jeroen van der Gun and and made other fixes and
improvements to PDF handling.

If you wonder why we don’t have browser plugin for Internet Explorer,
the explanation is simple: no-one has written the necessary code.

SumatraPDF has been brought to you by [SumatraPDF
developers](http://www.ohloh.net/p/4623/contributors).
