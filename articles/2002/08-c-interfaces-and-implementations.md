Id: 1361
Title: C Interfaces and Implementations
Tags: book,programming
Date: 2002-08-03T02:18:21-07:00
Format: Markdown
--------------
Finished reading [C Interfaces and Implementations](https://www.amazon.com/exec/obidos/ASIN/0201498413/).

Good, advanced book on C that teaches how to design and
implement good, re-usable interfaces.

It describes design and implementation of 24 interfaces (e.g. better memory
allocation primitives, better string-handling primitives, lists, dynamic
arrays, bit vectors, sequences, threads, etc.).

The book is not about analyzing and designing algorithms or data structures
but about practical issues of how to design good interface (which is
surprisingly hard).

It's very to the point. It doesn't have the flowery prose of many O'Reilly books.

Two tips from the book worth remembering:

-   the only way to get good performance is by profiling. Always prefer
    a clean code over what you think is more efficient code. First you
    might be wrong (as most people are) and second you might not need
    the speed in the first place (if you speed up by 100% code that only
    takes 1% of total execution time, you'll improve total execution
    time by only imperceptible 0.5%).
-   use `assert` generously

