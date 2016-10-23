Id: 45002
Title: Web server in C#
Tags: .net
Date: 2009-11-03T21:50:13-08:00
Format: Markdown
--------------
If you need a simple web server for .NET,
[this](http://www.codeplex.com/webserver) is a solid piece of software.

I’ve needed an internal web server for my C\# application and re-using
someone else’s code certainly beats writing it from scratch.

It’s not spelled out on the web page, but it comes in two flavors:

-   lite, with only http serving core. When you download the code, look
    for it in branches directory
-   standard, which includes additional capabilities like templating
    engine

When I added it to my project, I didn’t know about the lite version, so
I manually added the files to the project and removed files I didn’t
need. If I were starting from scratch, I would just use the lite
version.

So far I’m happy with the code - it does exactly what it advertises and
saved me a bunch of time I didn’t spend writing web server code.

It’s distributed under Apache 2.0 license so can be used in commercial
software as well.
