Id: 302
Title: Compile-time asserts in C
Tags: c,programming
Date: 2006-01-13T16:00:00-08:00
Format: Markdown
--------------
Asserts are good. Compile-time asserts are even  better. Here's how to define and use them in C.

```c
#ifndef CASSERT
#define CASSERT( exp, name ) typedef int dummy##name [ (exp ) ? 1 : -1 ];
#endif
CASSERT( sizeof(int16_t) == 2, int16_t_is_two_bytes )
```
