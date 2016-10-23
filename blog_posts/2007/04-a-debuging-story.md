Id: 970
Title: A debugging story
Tags: debugging,windbg
Date: 2007-04-29T14:38:15-07:00
Format: Markdown
--------------
or how a debugger sometimes works against you.

[SumatraPDF][1] was crashing on some PDF files and I didn't know why. On
Windows you can have just-in-time debugging i.e. tell the system to
automatically launch a debugger (Visual Studio or [WinDBG][2]) when a program
crashes but the clues I could get at the time of a crash (callstack and other
current state of the program) were not enough to figure out the cause
(especially given that it was in core PDF rendering code that I didn't write).

Usually I would just set a breakpoint just before the place of the crash and
work backwards from that. Unfortunately, when running under the debugger
(either Visual Studio or WinDBG) the crash didn't happen. The only good thing
about that was that it offered another clue: apparently something about the
system changes when the app is executed from within the debugger and it hides
the problem.

That stumped me for a while until I made a breakthrough: I figured out that if
I put DebugBreak() call in my app, it'll break into the debugger after the app
has started so I'll be able to debug.

With a few well-placed conditional breakpoints I was able to figure out that
the cause of the problem was uninitialized reference count on an object.

My theory is that when an app was executed from within a debugger, memory
allocator was using different flags and always zeroing malloc()ed memory
(which masked the refcount problem) while without the debugger it was random
data which exposed bad refcounting logic.

Lessons learned:

  * sometimes unexplained and very perplexing has an explanation
  * the debugger can work against you
  * [WinDBG][2] is awesome - it puts to shame anything I've used on Unix,
especially gdb. Learn how to use as it can come very handy when debugging hard
problems. Visual Studio is also a good debugger - I use both depending on the
kind of issue I'm debugging (for example WinDBG is faster to launch, has good
support for conditional breakpoints and an array of useful extensions)
  * refcounting is evil and don't let anyone tell you otherwise. You'll meet
people telling you that "refcounting makes programming easier". "Ha, ha, ha"
should be your response.

   [1]: http://blog.kowalczyk.info/software/sumatrapdf/

   [2]: http://www.microsoft.com/whdc/devtools/debugging/installx86.mspx


