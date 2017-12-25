Id: 1973
Title: memset() considered harmful
Tags: programming
Date: 2007-02-15T17:08:14-08:00
Format: Markdown
--------------
**Conclusion:** use `memzero(addr, size)` instead of `memset(addr, 0, size)`.

`memset` is an example of how blindly following established conventions leads
people to doing the wrong thing.

One of the most frequent uses for `memset` is to zero-out memory.

The problem with `memset` is that it's easy to swap the "value to set" with "count of values to set" arguments.

They're both ints, so the compiler won't complain. Humans are not good at remembering things of that nature.

It leads to many cases where `memset` is misused and ends-up being a no-op (setting 0 bytes of memory).

I've seen this in my own code, I've seen this in other people's code.

Solution is so simple: write a trivial wrapper `memzero(void *s,
int size)`.

It's a better name for the functionality and removes possibility of making this
particular mistake.

On Windows you don't even have to write it since it's there as `ZeroMemory()`.

Some unixes have `bzero()`, but it's a wierd name and the fucntion is not widely available so writing `memzero()` utility function is a good idea anyway.

