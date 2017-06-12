Id: 1923
Title: Code-name Monad and the value of different perspective
Date: 2005-10-26T12:19:17-07:00
Format: Markdown
--------------
ArsTechnica [wrote another excellent
review](http://arstechnica.com/guides/other/msh.ars/1 "wrote another excellent review"),
this time of upcoming Microsoft Windows Shell, code-name Monad.\
\
 It's a really cool shell, it blows every other shell in existence out
of the water. But this is not a link-blog, I provide commentary and what
Monad made me think about is the value of different perspective.\
\
 The first Unix shell [has been developed in
1971](http://en.wikipedia.org/wiki/Thompson_shell "has been developed in 1971"),
34 years ago. Bash, currently the most widely used shell on Unix was
[developed in
1987](http://en.wikipedia.org/wiki/Bash "developed in 1987"), a slight
improvement over [1977 Bourne
Shell](http://en.wikipedia.org/wiki/Bourne_shell "1977 Bourne Shell")
(thank god for whiskey and Wikipedia).\
\
 The most popular Unix shell, default on Linux, Mac OS X and most other
Unix versions, isn't much better than it was 30 years ago.\
\
 Open-source evangelists often argue that open-source is better than
closed source because openness brings experimentation, experimentation
brings improvements. They also argue that openness creates an
evolutionary way of writing software, where lots of changes in divergent
directions ultimately converge on the best solution, may the best piece
of code win. The problem with this view is:\

-   I don't want to wait a million years for a usable piece of software;
    I need it yesterday
-   evolution created humans so it must be great. But it also created
    pigs and, when you think about it, compared to humans everything
    else is a failure; I don't want software whose quality is comparable
    to the relative quality of a pig
-   you can only stretch a metaphor that much before it breaks. Software
    doesn't have DNA, what might work for nature doesn't have to work
    for software.

\
 The fact that Monad is so much better than bash is a slap in the
collaborative, evolutionary theory of software development.\
\
 If evolutionary way of writing software really is superior, then bash
should be a shining example. It's not an obscure utility no-one cares
about - it's an essential part of Unix.\
\
 It's not like no-one tried to improve the shell - we have Korn Shell, C
Shell, zsh, scsh (Scheme-based shell), Rsh, PWB Shell, Es shell,
Almquist shell, tcsh. Instead of advancing the default they mostly
succeeded in creating head-aches for system administrators.\
\
 During all this time all Microsoft had was cmd.exe, a sorry excuse for
a shell, a glorified DOS box.\
\
 So how come a few Microsoft programmers, in just few years,
leap-frogged 34 years of vigorous mutation and cross-breeding of Unix
shells?\
\
 I have a theory:\

-   the power of different perspective
-   the power of vision backed up by focused effort

\
 It's genuinely hard to come up with something better if there's a
strong history of good enough. It's hard to develop a different
perspective.\
\
 Unix shells are good enough and Unix heads for years derided Windows
for its clearly inadequate cmd.exe and claiming the superiority of Unix'
"lots of little tools, connected" philosophy. Coming from a position of
superiority, it's easy to become complacent and dismiss many flaws in
shell approach:\

-   number of shells in popular use causes more problems that it solves
-   shell is a really weak programming language
-   forking a processes to get a list of files is inefficient\
-   programs can only communicate via arbitrary, irregular text streams\
-   which is why we have hacks like awk. Please mommy, don't make me
    decipher awk scripts again, I'll be good, I promise
-   some simple and frequently done things like grepping through files
    matching a given pattern in all sub directories of a given directory
    often require syntax so complex that I don't even bother to remember
    what it is

\
 Apparently all Unix programmers prefer to adapt to those problems than
fix them once and for all. Microsoft, on the other hand, re-thought the
problem and fixed all those messy problems.\
\
 I might write about the power of vision backed up by focused effort
later.\
\
 And just in case you like to draw your own conclusions: please don't.
This is not "close source is better than open source" article. The world
is not black and white. There is some good open source and some good
closed software. There is shitty open source software and there is
shitty closed source software. If there's one grand conclusion to be
drawn here about open source is that it works, but for reasons different
than what the popular meme on the subject is.\
\

