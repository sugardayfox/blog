Id: 980
Title: Logging in WinDBG
Tags: debugging,windbg
Date: 2008-01-06T17:57:52-08:00
Format: Markdown
--------------
Logging is often useful when debugging software to help understand its
behavior.

Consider a problem I was having: Sumatra was crashing because of invalid
refcounting of fz_shade object. Refcounting took place only in two functions:
fz_keepshade() and fz_dropshade() so my plan of attack was to see who calls
those functions and from that hopefully figure out why there is a mismatch.

If I had a logging system with ability to dump callstacks already built into
the program, I could modify the program to add logging to those two functions,
re-run and look at the logs.

But I didn't so I used WinDBG as an ad-hoc logger.

The trick is to know the following:


  * .logopen $filename WinDBG command

  * .logclose WinDBG command

  * conditional breakpoints

Under WinDBG I set the following two breakpoints:

  * bp SumatraPDF!fz_dropshade "kb; g"

  * bp SumatraPDF!fz_keepshade "kb; g"

bp creates a new breakpoint on entry to a given function. A string "kb; g" is
composed of WinDBG commands executed when a breakpoint is hit. In this case I
just dump the callstack with kb and continue execution with g.

Then I used .logopen $filename so that everything WinDBG prints in output
window also gets written to a file and voila - an ad-hoc logging without the
need to modify and recompile the program.

And yes, after some staring at resulting log I fixed the crash.


