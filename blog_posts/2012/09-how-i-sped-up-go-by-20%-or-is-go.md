Id: 1404040
Title: How I sped up Go by 20% (or is Go really slower than Java?)
Tags: go,programming
Date: 2012-09-16T01:55:48-07:00
Format: Markdown
--------------
The story of a faulty benchmark
-------------------------------

I confess: I didn’t really speed up Go by 20%. However, if you were to
believe simplistic arguments, I did.

It all started when a comment on Hacker News claimed that Java is faster
than Go and pointed to [those
benchmarks](http://shootout.alioth.debian.org/u32/benchmark.php?test=all&lang=go&lang2=java)
as proof.

Indeed, every Go benchmark there is 2x to 10x slower than Java 7
program.

That seemed wrong to me - Go compiles to native code. Java, ultimately,
does too, but compilation happens at runtime, which is an overhead that
Go doesn’t have, so in principle Go should be faster.

I took a look at [Go
implementation](https://github.com/kjk/kjkpub/blob/master/gobench/revcomp.go).
I’m not an expert in implementing reverse-complement algorithm but I
immediately noticed that the code does a lot of unnecessary memory
allocations and copying.

I looked at [Java
implementation](http://shootout.alioth.debian.org/u64/program.php?test=revcomp&lang=java&id=6)
and what do you know: it’s implemented more efficiently, without
unnecessary memory allocations and copying. Plus it parallelizes the
computation.

I wrote [my own Go
implementation](https://github.com/kjk/kjkpub/blob/master/gobench/revcomp8c.go),
more similar to Java’s.

The result? On my Mac Pro I got about 20% speedup.

The argument of HN commenter was: look at those benchmarks, Go is 2x
slower than Java.

It’s a simplistic but unfortunately powerful argument: most people won’t
take the time to look at the code to determine that part of the issue is
that Java implementation is simply better than Go one.

By that simplistic argument I’ve improved speed of Go by 20% by writing
a slightly better implementation of the benchmark.

The story doesn’t end here. The above results were for single core x86
machine. If you look at results for [single-core 64-bit
machine](http://shootout.alioth.debian.org/u64/benchmark.php?test=all&lang=go&lang2=java),
Go actually wins on this and 2 other benchmarks.

The explanation is simple: 64-bit Go compiler is a little bit smarter
than 32-bit Go compiler.

Is Go really slower than Java?
------------------------------

There is no simplistic answer to this, but there are rules of thumb.

On equivalent programs, Go generates programs similar in speed or
faster. 64-bit version will be slightly better than 32-bit version.

Java garbage collector is much better than Go’s which explains why
allocation-heavy benchmarks (e.g. binary-trees) perform so much better
(there’s almost no computation there, just allocation of a lot of tree
nodes).

Those benchmarks aren’t currently representative of Go performance -
from cursory look it seems that C++ and Java implementations are
extremely optimized and Go implementations aren’t (I’ve submitted my
improved version, hopefully it’ll get incorporated).

Go is already competitive with Java and will only get better. Let’s not
forget that Java had over 20 years of investments in code generation and
garbage collection. Go is only 5 years old. There already were compiler
and garbage collection improvements since the latest released version.

One area where Go wins undeniably is memory usage - the programs use at
least order of magnitude less memory. Java is paying the cost of
sophisticated virtual machine.
