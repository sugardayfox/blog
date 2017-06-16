Id: 404003
Title: Writing a custom installer for Windows software
Tags: sumatra,programming
Date: 2011-02-06T14:43:27-08:00
Format: Markdown
--------------
A tale of 2 installers
----------------------

For version 1.3 of
[SumatraPDF](https://www.sumatrapdfreader.org/free-pdf-reader.html)
I’ve created a new installer.

The old installer:

![Old installer](//kjkpub.s3.amazonaws.com/blog/sumatra/sum-installer-old.png "Old installer")

New installer:

![New installer](//kjkpub.s3.amazonaws.com/blog/sumatra/sum-installer-new.png "New installer")

This post explains why I spent the time rewriting perfectly functional
installer.

How people write installers
---------------------------

If you write Windows software, you need to provide an installer. Almost
all installers are created by a tool (e.g.
[WiX](http://wix.sourceforge.net/), [Nsis](http://nsis.sourceforge.net),
[Inno Setup](http://www.jrsoftware.org/isinfo.php) and many others).

Apparently people think that writing an installer is so difficult that
it has to be simplified with a tool.

I’m here to advocate another way: write your own installers, people.

It’s not that hard.
-------------------

As everyone else, I did my installers with a tool.
[NSIS](http://nsis.sourceforge.net) was my weapon of choice: it’s free,
(relatively) easy to learn, (relatively well) documented and since other
people used it, there are examples to learn from.

I used it in
[SumatraPDF](https://www.sumatrapdfreader.org/free-pdf-reader.html),
[OpenDNS client](https://github.com/opendns/dynamicipupdate) and was
reasonably happy with it. However, it wasn’t all roses.

NSIS is great if you don’t try to do anything that it doesn’t support
out of the box. Once you stray outside the box, it’s hours spent trying
to find an existing NSIS code to do what you want and if you’re
unsuccessful, learning how to write code in NSIS’ awful macro language
or write extension dlls.

Two issues I ran into were:

-   since you can’t overwrite an executable file of a running program,
    it’s a good idea to kill running instance of the program in the
    installer, to avoid a failure in the middle of installation. There’s
    no built-in functionality to kill processes in NSIS. After spending
    way too much time trying to find an existing solution, I ended up
    writing an extension dll with necessary support command.
-   for my .NET app I had to install .NET 4.0 runtime if it wasn’t
    already installed. Doing that in C code is simple: check registry
    keys to determine if .NET runtime is installed and execute the
    installer if it isn’t. This is a common problem and I found several
    solutions on the net but each one was slightly different, since I’m
    not proficient in NSIS macro language they were not easy to tweak
    for my purpose and they were all outdated (e.g. showing how to
    install .NET 2.0 while I needed .NET 4.0).

I realized that an installer doesn’t do anything complicated: it shows a
window, writes a few registry keys a copies a few files. I figured that
writing installer in C++ from scratch wouldn’t take more time that I’ve
spent in total learning about NSIS.

I was right: it took me [9
days](http://code.google.com/p/sumatrapdf/source/list?path=/trunk/src/installer/Installer.cpp&start=2535)
to write shippable version of Sumatra’s custom installer. Most of it was
support code so writing future installers for my other software will be
even faster since I’ll reuse most of the code.

Why write custom installer?
---------------------------

I already had a functioning installer so why fix something that isn’t
broken? For me, the reasons where:

-   provide a better user experience
-   no frustration due to NSIS limitations

Installer is the first thing a potential user of your application sees.
What will be his first experience with your app?

Will it be terrible, like those installers that take half an hour to
unpack a few files?

Will it be
[boring](http://sethgodin.typepad.com/seths_blog/2009/01/youre-boring.html),
just like most other installers?

Or will it be
[remarkable](http://sethgodin.typepad.com/seths_blog/2003/06/remarkable_is_w.html),
unique, extra pleasant? An experience worth talking about to others?

Writer Joe Konrath [stress the
importance](http://jakonrath.blogspot.com/2010/12/holiday-ebook-buying-guide.html)
of a good book cover. An installer is the equivalent of a great book
cover.

Writing a custom installer doesn’t guarantee it’ll be exceptional, but
using NSIS guarantees that it’ll be boring, just like every other
installer created with NSIS.

In SumatraPDF’s case, I’ve made an usability improvement: installation
is a one-click affair (old installer required a few clicks).

On the whimsicality front, I’ve added a small startup animation, which
shows the program letters one by one. It’s only a small thing, I haven’t
spent much time on it, but the possibilities for interesting animations
are endless: scrolling text, 3d objects rotating, fireworks, star fields
etc. Old school [demos](http://www.pouet.net/) and
[processing](http://processing.org/) are a rich source of inspiration
for possible effects.

If you want to learn how a custom installer is done, the [source is out
there…](http://code.google.com/p/sumatrapdf/source/browse/#svn/trunk/src/installer)
