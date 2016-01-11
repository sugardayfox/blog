Id: 1286
Title: Sane #include hierarchy for C and C++
Tags: programming,c,c++
Date: 2007-06-24T17:00:00-07:00
Format: Markdown
--------------
In one large C/C++ project I worked on \#include files were so haphazard
and complex that sometimes rearranging the order of \#include statements
would cause compilation error that no-one could comprehend and fix in
any other way than reverting to previous order.

Every now and then an ambitious soul would set out to fix the problem
once and for all. And fail miserably shortly after.

The code in this project was quite good - it’s just maintaining good
include hierarchy requires a lot of effort and once you get into bad
state, its hard to get out of it.

The best weapon is being anal-retentive from the very beginning of your
project.

It’s important to reduce number of dependencies between header files,
especially circular dependencies.

WebKit project [spells the
rules](http://webkit.org/coding/coding-style.html):

-   All C/C++ files must \#include “config.h” first
-   All files must \#include the primary header second, just after
    “config.h”. So for example, Node.cpp should include Node.h first,
    before other files. This guarantees that each header’s completeness
    is tested, to make sure it can be compiled without requiring any
    other header files be included first.
-   Other \#include statements should be in sorted order (case
    sensitive, as done by the command-line sort tool or the Xcode sort
    selection command). Don’t bother to organize them in a logical
    order.

File `config.h`:
```
#ifndef COMMON_H_
#define COMMON_H_

/* This is a header included by every single file.
 It should include definitions of most basic types and utilities
 (like asserts and simple templates. It usually means most common
 header files for C standard library like <stdio.h> etc. and
 <windows.h> on Windows.
 Resist temptation to include too much in this file - it’ll lead
 to insanity.
*/
#endif
```

File `foo.h`:
```c
#ifndef FOO_H_
#define FOO_H_
/* To speedup compilation times avoid #include’ing header files.
 You can use any definitions in config.h. If possible, use forward
 declarations instead of including header file with definitions */

#endif
```

File `foo.c`:
```c
#include “config.h”
#include “foo.h”
/* rest of include files sorted alphabetically */
```
