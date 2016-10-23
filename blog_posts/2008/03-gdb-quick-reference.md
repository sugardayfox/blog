Id: 141
Title: gdb quick reference
Tags: debugging,gdb
Date: 2008-03-13T06:55:52-07:00
Format: Markdown
--------------
gdb quick reference.

```
attach $pid
gdb exe core
target remote 10.0.1.47:7777
set solib-search-prefix /dev/null
set solib-search-path /foo:/bar
gdbserver 0.0.0.0:7777 —attach $PID
info locals
info threads
info reg
disass

p/fmt expr
fmt: x - int in hex
 d - signed int
 u - unsigned int
 t - binary
 a - as addr, info symbol $addr
 c - char
 f - floating point
 s - string
e.g.: p/x $pc

x/nfu addr
 n - count
 f - format (x, d, u, t, a, c, f, s)
 u - unit size (b - byte, h (2b), w (4b), g (8b)
```
