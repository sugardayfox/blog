Id: 319001
Title: Why you shouldn't write Mac programs in QT
Date: 2010-08-22T15:46:25-07:00
Format: Markdown
--------------
[Toggl](http://www.toggl.com) is a time-tracking web application that I
just found. You can use it to measure how much time you spend on tasks.

They also have a desktop client for Mac, Windows and Linux, which is a
nice touch.

Except… the teeny-tiny desktop client, which has about 20k worth of
logic in it, is 114 MB big. It’s an obscene amount of code even in
today’s world, where neither programmers nor users care that much about
overhead. 106 MB of that is used by 9 Qt libraries.

The tragic thing is: with a little more (or even less) effort, it could
have been a tiny application. I understand the difficulty of delivering
a desktop app for 3 different OSes and the allure of using a framework
that promises to “write once deploy everywhere”, but Qt is not easy to
learn or easy to write code in (being a C++ monster that it is).

As far as I can tell Toggl’s desktop client is a window that embeds a
web browser and the UI is implemented in HTML with a bit of native logic
driven from it. Writing an application like that in C\# on Windows or
Cocoa on Mac is a weekend effort even by someone who doesn’t know those
platforms very well. The complexity of such program is at the level of
an introductory example for beginners.

Someone who knows what they’re doing would spend less time writing 2
apps using the best technologies on a given platform than it takes to
setup a build environment for Qt. If written in Cocoa or C\# they would
end up being tens of kilobytes in size, not hundreds of megabytes.
