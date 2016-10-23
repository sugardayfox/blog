Id: 1966
Title: On how I improved Sumatra performance by ~60%
Tags: sumatra,optimization,programming,profiling
Date: 2006-09-02T18:01:53-07:00
Format: Markdown
--------------
Improving performance of code you didn't write is fun. It's especially
fun if you can find a way to make a [small
change](https://bugs.freedesktop.org/show_bug.cgi?id=8112) and improve
performance by \~7%. It doesn't happen often, but when it does happen,
it feels good.\
\
[Sumatra](http://blog.kowalczyk.info/software/sumatrapdf/), my PDF
viewer for Windows, uses an existing open-source
[poppler](http://poppler.freedesktop.org/) library for most of its work
(parsing and rendering PDF files). Almost all execution time is spent in
poppler code. I didn't write the code yet I was able to improve some
specific scenario (loading of a PDF file) by about \~60% using rather
small and un-intrusive changes (i.e. without rewriting a lot of code).\
\
This post describes mechanics of performance improvement process I
used.\
\
Before you start improving performance, you need to decide what exactly
do you need to improve. You don't want to waste time [improving idle
loop](http://www-128.ibm.com/developerworks/power/library/pa-unrollav3/).\
\
In Sumatra case that was obvious. The things that people care about are:
how long does it take to load the PDF file and how long does it take to
render a given page.\
\
Another thing you need before even starting to improve the performance
is to have a way to reliably measure changes in performance. In my case
I had to write a test program that would take name of PDF file as input
and dump statistics about how much time, in milliseconds, did it take to
load the file and to render each page.\
\
In order to see how my changes affect performance I would build a
reference version of test program without any changes, build a version
with my changes (making sure I use the same compiler settings for both
versions), run them on the same PDF file and compare results.\
\
Getting results from a single run is not good enough. Given
multi-tasking nature of Windows, running time is only an approximation
of performance. Also, some parts of the code heavily depend on cache.
For example, loading PDF for the second time takes only a fraction of
time because the file is most likely cached in memory and memory access
is orders of magnitude faster than reading from disk.\
\
Therefore it's important to run the same test serveral (say 10) times
and calculate averages, rejecting values that fall too much outside of
average. Rather arbitrarily I decided to do a two-pass filtering when
calculating averages. In first pass I reject all values that differ more
than 45% from average and in second pass I reject all values that differ
more than 10% from the average. I don't have justification for those
values, they seem to work for me. The reason for two-pass filtering is
that one really outrages value in a short run might skew average so much
that all values differ from average a lot (e.g. \>20%).\
\
Even with multiple runs I see 1-1.5% changes when running the same
executable so I don't get excited when I see 1.5% improvement with my
changes - those fall into line noise.\
\
Another trap are CPUs with dynamically adjustable frequency (standard on
laptops since they save power). They completely ruin ability to use
execution time for comparison. Don't use laptops for running tests or
turn off this CPU feature.\
\
I used to run tests manually, copy data to spreadsheet and making
analysis there. That got old fast, so I wrote a python script to
automate the process. The script takes as arguments names of 2
executables to compare, name of PDF file, runs both tests programs on
this PDF multiple times (interleaved, to help even out possible CPU load
changes from other activity on the computer) and just gives me the
summary of results.\
\
Measuring execution times tells you how much work is being done but what
you need to know is which parts of the code are doing the most work.
Again, there's no point optimizing code that is not executed very often.
Speeding by 50% a piece of code that contributes to 1% of execution time
of your scenario only improves execution time by 0.5%. The best way to
get detailed performance information is to use a good profiler. On
Windows [AQTime](http://www.automatedqa.com/products/aqtime/) is quite
wonderful, Visual Studio Team (aka. Expensive) Edition has one. On Unix
I've read good things about [valgrind](http://valgrind.org/),
[oprofile](http://oprofile.sourceforge.net/)and
[dtrace](http://www.sun.com/bigadmin/content/dtrace/).\
\
The profiler will show you which functions take most of the time, how
many times they are called, call trees etc. This information is not only
useful to identify which code needs work but also can help you
understand how the code works. If code is algorithms + data then
software is code + execution paths.\
\
Then there is the hard (but most fun) part of figuring out how to change
the code to get speedups.\
\
There is occasional heartbreak. Some of my attempts at improvements only
gave me spectacular crashes (as always a result of [ignorance or
carelesness](/articles/sourceOfBugs.html)).\
\
On the bright side, it also happened that I had improvements ideas that
would require extensive changes but after staring more at the code and
profiler results, I found much smaller change with a similar speedup
potential.\
\
Summary of important points about performance optimization:\

-   optimize scenarios that matters to your users
-   write tools to make evaluating performance changes possible
-   write tools to automate performance evaluation
-   use profiler to see which code needs to be optimized
-   use profiler to understand how code works
-   use your brain to make changes to the code\
-   make sure your changes do improve performance. Flying blind is not
    good enough

If you're curious about what the specific changes for poppler were, you
can look at the patches I've submitted and their description: [\~19%
speedup](https://bugs.freedesktop.org/show_bug.cgi?id=7808), [\~25%
speedup](https://bugs.freedesktop.org/show_bug.cgi?id=7910), [\~2%
speedup](https://bugs.freedesktop.org/show_bug.cgi?id=8111), [\~7%
speedup](https://bugs.freedesktop.org/show_bug.cgi?id=8112).\
\

