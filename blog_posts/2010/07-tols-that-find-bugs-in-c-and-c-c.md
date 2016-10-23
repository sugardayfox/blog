Id: 266001
Title: Tools that find bugs in c and c++ code via static code analysis
Tags: programming,c,c++
Date: 2010-07-07T21:18:43-07:00
Format: Markdown
--------------
I love tools that help me write better code. C and C++ are especially
hard on humans and provide the biggest amount of rope to hang yourself
with so tools that can point problems in C/C++ code are especially
useful. Static code analysis tools can analyze C/C++ code without
running it (hence the “static” part) and find a multitude of issues,
like memory leaks, using uninitialized variables, freeing memory twice
etc.

Here are the tools that I recommend because I’ve used them and they
worked.

cppcheck
--------

I’ve learned about
[cppcheck](http://sourceforge.net/apps/mediawiki/cppcheck/index.php?title=Main_Page)
just recently because someone [opened a
bug](http://code.google.com/p/sumatrapdf/issues/detail?id=984) against
Sumatra with a result of running a cppcheck scan. It’s free, open source
and even has a GUI for Windows (it could be improved, though).

It’s very easy to use: you just give it a file name or a directory, it
runs its magic on the C/C++ files and tells you what the problem is and
file name/line number where it happens.

It tells you not only about incorrect code (like memory leaks or using
uninitialized variables) but also about style issues (e.g. it will
suggest making a C++ class member a const function if it can be const
etc.).

I found error messages to be clear. Cppcheck complains about things that
I don’t always consider problematic enough to fix, but it’s easy to just
ignore them.

Since it’s (also) a cmd-line tool, it could be integrated with a build
system.

A unique property of cppcheck is that it can be run on stand-alone
files. The advantage is that it’s really easy to use (you can, for
example, check unix or mac source code on windows (and vice-versa)). The
disadvantage is that unlike other tools described here it probably can’t
do as good of a job analyzing things since it has less information to
use.

Given the price (free) and how easy it is to use, there’s no excuse for
a C/C++ programmer to not run it from time to time over their code
bases.

Clang Static Analyzer
---------------------

[Clang Static Analyzer](http://clang-analyzer.llvm.org/) is also free
and open-source. It’s part of Apple’s [clang](http://clang.llvm.org/)
project.

Latest builds of Xcode (3.2 and later) have some version of CSA
integrated so running it from inside XCode is super easy. After the run
it annotates the source code with issues found.

It can be also run as a [stand-alone
tool](http://clang-analyzer.llvm.org/scan-build.html) which “hijacks”
compiler invokations from a make or xcodebuild tools. This mode of
operations has an advantage of being easy to use if your project is
using makefiles (or Xcode) and that it compiles exactly the same code
with the same flags as real compiler.

The disadvantage is that if you have a different build system, you can’t
use it.

Stand alone tool produces html file which shows source code annotated
with information about found issues.

For mac os pre-build binary of stand-alone tool is available from [the
main page](http://clang-analyzer.llvm.org/scan-build.html) (and is most
likely newer than the copy that comes with Xcode). For other platforms
it can be built from source.

Visual Studio
-------------

Premium and Ultimate (i.e. expensive) versions of Visual Studio come
with static code analysis built in. The technology comes from previously
stand-alone tool called
[PREfast](http://blogs.msdn.com/b/oldnewthing/archive/2010/06/15/10024931.aspx).

To use it you you need to provide
[/analyze](http://msdn.microsoft.com/en-us/library/ms173498.aspx) option
to the compiler or configure it in UI via target settings properties
(under Code Analysis).

Then you perform a build and Visual Studio tells you about problems it
found.

Here’s a [more extensive
report](http://www.cs.auckland.ac.nz/~pgut001/pubs/sal.html) on using
Visual Studio’s static analyzer.

Coverity
--------

[Coverity](http://www.coverity.com/products/static-analysis.html) is one
of the first static analyzers for C/C++ code. I used it at Palm and it
works well. The problem is that it’s probably expensive (Coverity
doesn’t even publish the prices on the website, you have to go through a
sales person, which is never a good sign).

Also, I don’t know how easy it is to use. At Palm the integration with
build system was done by the build team and I only interacted with the
tool through (quite nice) web based ui.

How do they compare
-------------------

I don’t know because I haven’t run a comparative analysis.

What I do know is that they all found issues with the code I was working
on so they’re valuable.

Given that cppcheck and Clang Static Analyzer are free and easy to use,
they should be used by all C/C++ programmers, ideally integrated into
the build process so that they’re easy to run. Preferably they should be
run automatically after each checkin.
