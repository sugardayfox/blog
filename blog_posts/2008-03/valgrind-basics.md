Id: 1043
Title: valgrind basics
Tags: unix,debugging,programming
Date: 2008-03-13T17:19:08-07:00
Format: Markdown
--------------
[Valgrind](http://valgrind.org) basics is a tool that can find bugs in
compiled code by running compiled executable and instrumenting it.
According to docs works best on code compiled with `-O0` or `-O1` and
might report false positives with highly optimized code (i.e. compiled
with `-O2`).

**Basic usage**
```
valgrind \${prog-to-run} ${arg1} ${arg2} …
```

**Checking for memory access errors and memory leaks**:

```
valgrind —tool=memcheck —leak-check=full —show-reachable=yes
—num-callers=16 ${prog-to-run}
```

**Profiling memory usage:**

```
valgrind —tool=massif ${prog-to-run}
gv massif.${pid}.ps
most massif.${pid}.txt
```

**Profiling CPU usage:**

```
valgrind —tool=callgrind ${prog-to-run}
```
or:

```
valgrind —tool=callgrind —instr-atstart=no ${prog-to-run}
```

and then:

```
callgrind_control -i on
callgrind_control -i off
```

It generates `callgrind.out.<pid>` file (where `<pid>` is process id of
the program). Use `kcachegrind` GUI tool to view the results.

**Profiling sections of the code**

-   `#include <valgrind/callgrind.h>`
-   wrap the section to profile with `CALLGRIND_TOGGLE_COLLECT` macros
-   add `--collect-atstart=no` so that valgrind doesn’t profile from the
    beginning but only the sections of the code within
    `CALLGRIND_TOGGLE_COLLECT`

In other words, if your code looks like:
```
#include <valgrind/callgrind.h>

CALLGRIND\_TOGGLE\_COLLECT
foo();
CALLGRIND\_TOGGLE\_COLLECT
```
your program is profiled only when foo() function is running.

There’s also an option `--instr-atstart=no` and
`CALLGRIND_START_INSTRUMENTATION` and `CALLGRIND_STOP_INSTRUMENTATION`
macros. In theory you can use this option and insert those 2 macros into
the code to speedup profiling (instrumentation slows down valgrind) but
in practice valgrind 3.2 crashed when I was using this option.
