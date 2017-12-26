Id: 952
Title: Order of #include headers in C/C++
Tags: programming, c++
Date: 2006-08-15T09:23:28-07:00
Status: notimportant
Format: Markdown
--------------
One thing I've learned: maintaining good #include hierarchy requires eternal vigilance.

It's easy to slack off and end up with a mess like:
* circular dependencies (you must include both files in order to compile any of them)
* transitive inclusion. Files that compile only because some other file happened to have been included somewhere in #include chain.

This mess is not a theoretical problem: it bites you when you modify the code and suddenly it doesn't compile because of bad of #include dependencies.

Compilation problems caused by messed up dependencies are hard to figure out.

One big project I've worked on had this problem and a running joke was
that every couple of months some developer would get determined to fix
it once and for all by cleaning up headers.

After all, how hard can it be? Turns out it was very hard and no-one succeeded.

For that reason I cringe every time I see `#include <stdafx.cpp>` - it's
a free ticket to future dependency hell.

A trick I recently settled upon helps to keep clean \#include hierarchy.
In the past (for no reason I can remember) I would put \#include for
system includes (like or ) first in my \*.c files. Those days the first
\#include in module foo.c is for "foo.h".

Why?

The golden rule for #include files is that if a module bar.c uses
foo.c, everything needed to compile foo.c should be defined in foo.h.
Chances are that `foo.h` uses definitions defined in system includes. If
all places that `#include "foo.h"` also include those system includes
before foo.h, things will compile just fine but only by accident.

Which is not a problem until you forget to #include those system
includes and are faced with weird ("it used to work just fine") compiler
errors.

Including it's own \#define as the first thing helps to spot those
mistakes early.
