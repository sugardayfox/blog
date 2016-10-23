Id: 50001
Title: Drobo Dashboard and mysterious mac slowdowns
Tags: mac
Date: 2009-11-24T20:41:47-08:00
Format: Markdown
--------------
I have a mac pro (early 2008 model) running OS X 10.6 and recently it
has been plagued with mysterious slowdowns: after a couple of hours it
would develop a huge latency when responding to mouse clicks.

For example it would take literally seconds to show menu after clicking
on a menu bar or switch to an application after clicking its window. Not
an experience I would expect from 8 core machine with 12GB of RAM.

A reboot would fix it but only for a few hours. After several days of
frustration I figured out the problem: Drobo Dashboard 1.6.1 was the
culprit.

There are two solutions.

One is to kill DDServiced process when the slowdown starts happening
(`sudo killall -9 DDServiced`). It’ll restart itself so things will keep
working and the slowdown will be gone for several hours.

Another one is to uninstall Drobo Dashboard (at least until Drobo fixes
their software) following [those
instructions](http://support.datarobotics.com/app/answers/detail/a_id/107).
Drobo runs just fine without the dashboard and its usefulness is lost on
me.

I found a few [twitter mentions](http://twitter.com/search?q=DDServiceD)
about this problem. According to [this
tweet](http://twitter.com/dionon/statuses/6000094679) upcoming 1.6.6
release will fix this problem.
