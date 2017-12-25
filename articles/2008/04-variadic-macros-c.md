Id: 1253
Title: Variadic Macros (C++)
Tags: visual studio,c++
Date: 2008-04-04T16:36:43-07:00
Format: Markdown
--------------

Variadic macros are function-like macros that contain a variable number
of arguments.

To use variadic macros, the ellipsis may be specified as the final
formal argument in a macro definition, and the replacement identifier
**\_\_VA\_ARGS\_\_** may be used in the definition to insert the extra
arguments. **\_\_VA\_ARGS\_\_** is replaced by all of the arguments that
match the ellipsis, including commas between them.

The C Standard specifies that at least one argument must be passed to
the ellipsis, to ensure that the macro does not resolve to an expression
with a trailing comma. The Visual C++ implementation will suppress a
trailing comma if no arguments are passed to the ellipsis.

Support for variadic macros was introduced in Visual C++ 2005.


#### Example:

```c++
    // variadic_macros.cpp
    #include
    #define EMPTY

    #define CHECK1(x, ...) if (!(x)) { printf(__VA_ARGS__); }
    #define CHECK2(x, ...) if ((x)) { printf(__VA_ARGS__); }
    #define CHECK3(...) { printf(__VA_ARGS__); }
    #define MACRO(s, ...) printf(s, __VA_ARGS__)

    int main() {
       CHECK1(0, "here %s %s %s", "are", "some", "varargs1(1)\n");
       CHECK1(1, "here %s %s %s", "are", "some", "varargs1(2)\n");   // won't print

       CHECK2(0, "here %s %s %s", "are", "some", "varargs2(3)\n");   // won't print
       CHECK2(1, "here %s %s %s", "are", "some", "varargs2(4)\n");

       // always invokes printf in the macro
       CHECK3("here %s %s %s", "are", "some", "varargs3(5)\n");

       MACRO("hello, world\n");
       // MACRO("error\n", EMPTY);   would cause C2059
    }
```

#### Output

```
    here are some varargs1(1)
    here are some varargs2(4)
    here are some varargs3(5)
    hello, world
```
