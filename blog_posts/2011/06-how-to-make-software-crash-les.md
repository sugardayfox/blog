Id: 566003
Title: How to make software crash less
Tags: programming
Date: 2011-06-30T15:32:59-07:00
Format: Markdown
--------------
You don’t want your software to crash, do you? This post describes my
experiences in making desktop Mac and Windows software crash less.

Know thy crashes
----------------

The most important step to fixing your crashes is to know about them.

Given how complex and varied our desktop operating systems are (3 major
version of Windows in active use, thousands of minor and major ways that
each installation of Windows or Mac can be different) our ability to
comprehensively test software is not good.

Sure, if you’re Microsoft or Adobe you can reinvest some of the revenue
to hire an army or testers, setup compatibility labs etc. but for a
small company or a single developer this is not realistic.

Bugs most often lurk in untested code and even a very good testing
effort is unlikely to encounter all the many things that can go wrong in
real life.

Get the crash reports automatically
-----------------------------------

Very few people bother to tell software vendors about crashes. They just
shrug and restart. The only realistic way to be informed about crashes
is to automatically gather crash reports without user involvement.

This is a proven idea. Microsoft was one of the pioneers of using this
technique for Windows OS and they publicly praise it for letting them
fix the most frequent crashes and increase stability of Windows. Many
clued in organization do it as well: Apple, Mozilla, Google.

Getting the crashes - the mechanics
-----------------------------------

Regardless of the platform, the solution involves two parts:

-   a server, which accepts crash reports from the software
-   code inside the software itself. It gets activated when a crash
    happens and sends crash report to the server

### The server

The server part is simple. You can use any technology to write it, any
protocol you want.

Personally I use [App Engine](https://appengine.google.com/), re-use
HTTP POST protocol and run on standard HTTP port (80) for maximum
compatibility with client-side firewall software.

It’s literally few lines of code to parse incoming POST requests, store
them in the database and provide a basic web interface for easy
browsing. As an additional bonus, App Engine is free if your traffic is
small enough,

If you write in PHP and run on a small VPS, it’ll work just as well.

### Client side on Windows in C\# apps

In C\# I setup a global exception handler (in WPF it means setting up
handlers for App.DispatcherUnhandledException and
AppDomain.CurrentDomain.UnhandledException).

I then use HttpWebRequest class to send exception message and callstack
as HTTP POST (using multipart/form-data).

The code is less than hundred lines and took mere hours to write.

### Client side on Windows in C++ apps

In C++ everything is substantially harder. The big picture is similar:

-   install global handler for unhandled exceptions (with
    SetUnhandledExceptionFilter())
-   in that handler generate the crash report and HTTP submit it to the
    server

The details are substantially more complicated. The code is more than a
thousand lines and took me days to perfect.

The biggest issue was that, unlike in C\#, you don’t get readable
callstack in native code. You need symbols (.pdb files) for that. My
solution is convoluted, but works:

-   during build process, I archive .pdb files on the server (I use S3,
    but any web server will do)
-   in the crash handler I download the symbols locally, on demand
-   I create crash report containing human-readable callstacks of all
    threads and some other useful information
-   I HTTP POST it to the server

You can re-use my work. The code is part of my
[SumatraPDF](http://blog.kowalczyk.info/software/sumatrapdf/)
open-source project. Most of it is in
[CrashHandler.cpp](http://code.google.com/p/sumatrapdf/source/browse/trunk/src/CrashHandler.cpp)
and, unlike most of Sumatra, is under liberal BSD license.

### Client side on Mac OS X

The good thing about Mac is that it already creates human-readable crash
reports for you. They are stored in \~/Library/Logs/CrashReporter/
directory.

There’s no need to handle crashes yourself. I just check at startup if
there’s a new crash report for my app (the files are named after your
application) from a previous run. If there is, I submit it to the server
(and delete so that it’s not sent multiple times).

The alternatives
----------------

While the general idea is always the same, there are different ways of
implementing it.

On Windows a simpler solution is to capture so-called minidumps (using
MiniDumpWriteDumpProc() Windows API) instead of going to the trouble of
generating human-readable crash reports client side.

I did that too. The problem with that approach is that you have to
inspect each crash dump manually in the debugger (e.g. WinDBG). I wrote
a python script that automated the process (you can script it by
launching cdb debugger with the right parameters and making it run
!analyze -v)).

Unfortunately, cdb is buggy and was hanging on some dump files. It’s
probably possible to work around with a timeout in the python script,
but at that point I stopped caring.

Windows provides native support for minidumps. Google took minidump
design and provided cross-platform implementation for Windows, Mac and
Linux, as part of [breakpad](http://code.google.com/p/google-breakpad/)
project.

Breakpad is the crash reporting system used by Google for Chrome and
Mozilla for Firefox. It contains both client and server parts for native
(C/C++ or Objective C) code.

I used it once for a Mac app. For Objective C I prefer the approach
described above as it’s simpler to implement, but I’m sure that’s a
solid and well tested approach.

On Windows, crash reports from your app are already sent to Microsoft as
part of [Windows Error
Reporting](http://en.wikipedia.org/wiki/Windows_Error_Reporting).
Apparently, it’s possible to for third party developers to get access to
those reports but I never did that, so don’t know what’s involved in the
process.

The SumatraPDF experience
-------------------------

So how well does it work in practice?

We’ve implemented the system described here in Sumatra 1.5. Sumatra is a
rather complicated piece of C++ code and quite popular (several thousand
of downloads per day).

Before 1.5 we had a system where we would save the minidump to a disk
and after a crash we would ask the user to report it in our bug tracker
and attach minidump to the bug report.

It became obvious that almost no one did that. We’ve only gotten few
crash reports from users in few months. Our automated system was sending
us tens of crash reports per day.

Once we knew about the problems, we could attempt to fix them. Some
problems we could fix just by looking at crash report. Some required
writing stress tests to make them easier to reproduce locally. Some of
them we can’t fix (e.g. because they are caused by buggy printer drivers
or other software that injects buggy dlls into our process).

We do know that we fixed some of the bugs. We can see that a new release
generates less crashes and by looking at crash reports we can tell that
some crashes that happened frequently in previous releases do not happen
anymore.

Building automated crash reporting system was the best investment we
could have made for improving reliability of SumatraPDF.
