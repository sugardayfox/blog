Id: 338001
Title: 8 habits for becoming a better programmer
Tags: programming
Date: 2010-11-05T18:21:53-07:00
Format: Markdown
--------------
This post is inspired by [a
question](http://news.ycombinator.com/item?id=1674103) “What habits made
you a better programmer”.

All the habits described below stem from one realization: I’m too stupid
to write bug free code so I need ways to compensate for my inevitable
human fallibility.

Code reviews
------------

I welcome code reviews because a second pair of eyes might spot mistakes
I’ve made.

Some people can’t stand their code being criticized. I avoid bruised ego
by assuming from the start that I’m too stupid to write correct code.

Using good tools
----------------

I seek out and use tools that help me find bugs automatically and tools
that help me understand my code better. Those tools include:

-   static code checkers like [clang
    analyzer](http://clang-analyzer.llvm.org/) or
    [cppcheck](http://sourceforge.net/apps/mediawiki/cppcheck/index.php?title=Main_Page)
    or [pychecker](http://pychecker.sourceforge.net/)
-   [Valgrind](http://www.valgrind.org/)
-   memory and cpu profilers
-   debuggers
-   [Source Insight](http://www.sourceinsight.com/) (a text editor)

Automated testing and continuous builds
---------------------------------------

Continuous build quickly alerts to mistakes that break the build.

Automated tests (unit tests, system tests) increase my confidence in the
correctness of the code and catch mistakes that cause regressions.

Stepping through all new code in the debugger
---------------------------------------------

I see a lot of macho anti-debugger posturing. The reasoning is: good
debuggers make it too easy to fix problems which makes you a sloppier
programmer. But this reasoning applies to unit tests as well (although I
don’t see anyone criticizing unit tests for that reason).

Stepping through newly written code in the debugger to double-check that
it behaves the way I expect is just another way to compensate for the
inevitability of writing buggy code.

Avoiding complexity
-------------------

Unnecessary, productivity sapping complexity is something that no one
would defend but popularity of C++ or boost shows that a definition of
unnecessary complexity is different for different people.

My bar for calling something complex is lower that many. I stay away
from complexity, both self-inflicted (like trying to be too clever when
implementing something, using multithreading when it’s not absolutely
necessary) or inflicted by the tool (e.g. advanced features of C++).

Diagnostic code built into apps
-------------------------------

I add diagnostics to my code. Logging, asserts in debug builds, crash
dump submission to my site for analyzing crashes that happen in the
wild. They all help to figure out the inevitable problems that you don’t
see in testing on your machine but happen on user’s computers due to
different configurations.

Writing readable code
---------------------

No one has writing unreadable code as a goal but it does happen because
writing readable code is harder and requires more care and attention
that just writing some code that works.

Readable code is important because I know that bugs will happen despite
my best efforts to prevent it. To fix them I will have to read my own
code long after I wrote it. Therefore I try to make my code as readable
as possible for my future self. The specific techniques involve:

-   balanced commenting. Not too much (so that it doesn’t detract from
    the code itself) but also not too little so that non-obvious
    decisions that are not captured in the code are explained.
-   no cryptic names for variables or functions
-   no long functions with complex logic
-   taking the time to make the code look consistent

Re-use high-quality code
------------------------

It’s much better if other people sweat writing code and fixing bugs. I
look for high quality, reputable code and use it whenever I can e.g. I
will use [SQLite](http://www.sqlite.org/) rather than writing my own
persistence layer. 
