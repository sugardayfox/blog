Id: 1937
Title: Pawn: yet another embedable language
Tags: programming
Date: 2006-01-14T15:05:00-08:00
Format: Markdown
--------------
I thought I knew about them all, but no, I just found another one: [Pawn][1].

Pawn is C-like language intended as an embeddable scripting language for
writing extensions for your C code. Not that we need yet another one.

I've only skimmed the website, but it looks fairly interesting. It's being
maintained, it comes with zlib license (i.e. open-source and can be used in
commercial apps), compiles to p-code (i.e. virtual machine), has interpreter
for p-code, optimized interpreter written in assembly as well as JIT
interpreter. It works under Windows and Unix. I would check it out (as an
alternative to lua or python) if I was writing an app that needs an embeddable
scripting language.

   [1]: http://www.compuphase.com/pawn/pawn.htm (Pawn)


