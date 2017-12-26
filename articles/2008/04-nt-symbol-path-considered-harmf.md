Id: 995
Title: _NT_SYMBOL_PATH considered harmful
Tags: windbg,debugging
Date: 2008-04-17T19:31:59-07:00
Format: Markdown
Status: hidden
--------------

When you debug your software on Windows, it’s sometimes useful to also
have access to symbols for Windows DLL. Microsoft is kind enough to provide
them for download and the easiest way to make them available locally is to
set `_NT_SYMBOL_PATH` environment variable so that it downloads symbol files
locally from Microsoft servers when needed. The right value is
`srv*c:\symbols*http://msdl.microsoft.com/download/symbols` (you can
change `c:\symbols` to any directory on your hard-drive).

I was running such setup for some time but observed that starting up a
debugging session too a long time. It was annoying and after a while it
hit me that I’ve caused this: even the simplest Windows program implicitly
links a lot of DLLs. Both Visual Studio and WinDBG feel the need to load up
symbols for all those system DLLs and it turns out it takes a lot of time.

I unset `_NT_SYMBOL_PATH` and I’m a happier person for it. Most of the
time I don’t need system symbols since the bugs are in my code, but sometimes
having them is useful (e.g. when the crash happens inside Windows). Fortunately
you can still have access them. When needed you can use `.sympat`h to set
the path (e.g.
`".sympath srv*c:\symbols*http://msdl.microsoft.com/download/symbols"`)
and then do `.reload` to force loading of symbols.
