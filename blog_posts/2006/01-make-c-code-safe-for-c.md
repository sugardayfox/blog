Id: 296
Title: Make C code safe for C++
Tags: c,c++,programming
Date: 2006-01-25T16:00:00-08:00
Format: Markdown
--------------
When C or C++ compiler compiles a function, it generates an object file that
contains generated assembly code identified by a symbol generated from
function name. Generating this symbol is called name mangling. C compilers
use very simple name-mangling: they prepend single "\_" to function name i.e.
function `int foo(int a)` is identified in object file by a symbol `_foo`.

C++ name mangling are much more complicated since C++ has to support classes,
overloaded functions etc.

It's also incompatible with C compiler's name mangling. It can be a problem.

If a header file for functions compiled by C compiler is included in a file
compiled by C++ compiler, the linker won't be able to link those object
files since the names of functions won't match due to different name mangling.

There is a solution to this problem.

First thing to know is that C++ compilers
define __cplusplus preprocessor symbol, so it's possible to use `#ifdef __cplusplus`
statement to compile parts of the code only by C++ compiler.

Second thing to know is that `extern "C" { /* block of code */ }` tells C++ compiler
to interpret a given block of code as C compiler would i.e. use C compiler's name
mangling for those functions.

All this boring explanation means that in order to make your C code safe for C++,
you need to wrap the code in header file with those statements:

```c++
#ifdef __cplusplus
extern "C"
{
#endif

/* ... your C code ... */

#ifdef __cplusplus
}
#endif
```
