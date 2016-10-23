Id: 256002
Title: Comparing program versions (in C# and Python)
Tags: programming,c#,python
Date: 2010-07-25T17:48:00-07:00
Format: Markdown
--------------
When you write an auto-update system for desktop software, you need code
that compares two version numbers (given as strings) and tells you which
version is greater.

This would be simple if we could stick with a simple version naming
scheme e.g. we could say that a version number is just an integer (1, 2
etc.) in which case it would be as simple as converting a string to an
integer and comparing those integers.

Version numbers in real life are more complicated. They can be
multi-part (e.g. "1.2.32). We have notions of alpha and beta builds etc.
Writing code to compare such program version numbers seem daunting.

There is, however, a relatively simple algorithm that supports
multi-level version numbers and notions of alpha, beta and rc builds.
Among other things, it knows the following relations between program
numbers:

-   1.1 \> 1.0
-   1.0.38 \> 1.0.5
-   1.1 \> 1.1rc \> 1.1b \> 1.1a (‘a’ stands for alpha builds, ‘b’ for
    beta and ‘rc’ for, well, rc builds)

The algorithm is quite clever but instead of describing it, I’ll just
include the Python and C\# implementations. I wrote them myself and I
put the code into public domain, so you can use it any way you want. I
use C\# implementation in my desktop software.

Go to:

-   [Python implementation](http://gist.github.com/458188)
-   [C\# implementation](http://gist.github.com/458186)

