Id: 1243
Title: variadic macros in msvc
Tags: programming,c++,msvc
Date: 2008-04-04T16:38:52-07:00
Format: Markdown
--------------
<div>

\
\>* Is there any way to support variadic macros in Visual Studio\
*\>* like gcc supports?  If not, how do people write stuff like\
*\>* custom asserts that take variable parameters?\
*\>\
\>* <http://gcc.gnu.org/onlinedocs/cpp/Variadic-Macros.html>\
*\>\
\
Depending on what you're after, you might find that the \_\_noop\
intrinsic will help you out.\
\
 \#ifdef \_DEBUG\
void \_MyDebugPrintf(const char\* fmt, ...);\
\#define DPF \_MyDebugPrintf\
\#else\
\#define DPF \_\_noop\
\#endif\
\
// later, ...\
\
        DPF("%s %s this works", FuncA(), FuncB());\
\
FuncA and FuncB (or parameters in general) are not evaluated\
in non \_DEBUG builds. The compiler ignores the statement.\
\
 This is obviously not variadic macros, but you can do quite a\
bit with it.\

</div>

<div align="right">

</div>
