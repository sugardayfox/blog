Id: 11
Title: SumatraPDF 3.0 released
Tags: releasenotes,sumatra
Date: 2014-10-29T20:56:38-07:00
Format: Markdown
--------------

We, [the SumatraPDF developers](http://www.ohloh.net/p/4623/contributors?sort=latest_commit&time_span=30+days) have released a version 3.0 of Sumatra, a multi-format reader (PDF, epub and mobi ebooks, comic books, etc.) for Windows.

You can download it from [official SumatraPDF website](https://www.sumatrapdfreader.org/free-pdf-reader.html)

The biggest change in this version is addition of tabs, contributed by Stefan Stefanov.

If you don’t like tabs, you can go back to the old UI using Settings/Options... menu

We added support for table of contents and links in ebook UI.

We added support for PalmDoc ebooks.

Comic books now support CB7 and CBT format (in addition to CBZ and CBR).

We added support for LZMA and PPMd compression in CBZ comic books

You can now save Comic Book files as PDF.

We swapped keybindings:

 * F11 : fullscreen mode (still also Ctrl+Shift+L)
 * F5 : presentation mode (also Shift+F11, still also Ctrl+L)

We added a document measurement UI. Press 'm' to start. Keep pressing 'm' to change measurement units.

We added new [advanced settings](https://www.sumatrapdfreader.org/settings.html): FullPathInTitle, UseSysColors (no longer exposed through the Options dialog), UseTabs

We replaced non-free UnRAR with a free RAR extraction library. If some CBR files fail to open for you, download unrar.dll from [rar website](http://www.rarlab.com/rar_add.htm) and place it alongside SumatraPDF.exe

We deprecated browser plugin. We don’t remove it if it was installed in previous version but both Chrome and FireFox are removing support for plugins so there’s no point in keeping it.

Finally, some of you really didn’t like the yellow background color. You’ve won: it’s now gray.
