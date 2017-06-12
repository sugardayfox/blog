Id: 1989
Title: gflags - a debugging story
Tags: programming,debugging,win32
Date: 2008-04-07T00:28:39-07:00
Format: Markdown
--------------
I hate crashes that disappear when run under the debugger and I had one
when porting mupdf to Windows.

It helps to know that there’s at least one reason for a changed behavior
under the debugger: it automatically triggers using debugging heap. While
debugging heap usually helps find problems, sometimes it does the opposite by
changing the details of memory allocation.

One helpful tool when debugging memory problems on Windows is gflags
which can enable page heap instrumentation for a given program. It works by
putting each allocation into a separate region of memory and putting a non-readable
page right after that. Also, upon freeing it makes the memory unreadable.
That way an overwrite of memory block while it’s still being used or accessing
the memory after it was freed will cause immediate crash.

The downside is that using gflags uses much more memory. But in those
days of cheap gigabytes it’s not a problem that can’t be solved with a couple
hundred bucks.

Basic usage of gflags.exe is simple: `gflags /p /full /enable foo.exe`

From now on foo.exe will always be run with this instrumentation turned
on. To disable, do `gflags /p /disable foo.exe`

To see which programs have page heap enabled, do `gflags /p`. gflags
offers many other option and you can learn about them via `gflags /?`. If you run
gflags without any options, you’ll get a (very confusing) GUI.

It worked like a charm. I got a crash on accessing freed memory and all
I had to do was to backtrack to where this memory was allocated to figure out
the problem.
