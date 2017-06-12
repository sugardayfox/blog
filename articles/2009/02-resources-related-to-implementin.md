Id: 8003
Title: Resources related to implementing programming languages
Tags: prog lang,programming
Date: 2009-02-21T14:38:25-08:00
Format: Markdown
--------------
http://ece-www.colorado.edu/\~siek/ecen4553/ - a course teaching
compiling python. Has assignments and links to resources (like info
about x86 etc.)

http://www.tinypy.org/ - small python implementation. register-based vm
inspired by lua, own gc

http://code.macournoyer.com/tinyrb/ - small ruby implementation, boehm
gc

http://code.google.com/p/v8/ - v8, fast JavaScript implementation,
compiling directly to native code from parse trees, own gc, C++,
embeddable, New BSD license

http://piumarta.com/software/lysp/ - really small lisp implementation,
boehm gc or its own, MIT license

http://github.com/why/potion/tree/master - potion, jit compilation, vm
based on lua, used ‘id’ object model, gc is reference counting, MIT
license

http://www.complang.org/ragel/ - compiles FSM from regexpes, can be used
for lexical parsing, used by tinyrb

http://www.hwaci.com/sw/lemon/ - lalr(1) parser generator, like yacc,
used by tinyrb
