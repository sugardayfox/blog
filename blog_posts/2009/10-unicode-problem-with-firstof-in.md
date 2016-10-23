Id: 41001
Title: Unicode problem with firstof in appengine/Django
Tags: appengine
Date: 2009-10-21T17:30:33-07:00
Format: Markdown
--------------
I spent a couple of hours tracking and fixing this problem.

A user reported an error in my
[fofou](http://blog.kowalczyk.info/software/fofou/) forum software for
appengine: a forum that had german text in tagline couldn’t be edited.
All appengine told me was this mostly useless error message:
“UnicodeEncodeError: ‘ascii’ codec can’t encode character u’\\xdc’ in
position 0: ordinal not in range(128)” coming from template.render.

By trial and error I narrowed down the culprit to this part of the
template: “{% firstof forum.tagline ”Tagline." %}". forum.tagline comes
from a datastore, so it’s unicode. Error message tells me that (for
whatever reason) template rendering is trying to convert that unicode
string to ascii using ascii codec and (rightfully) failing.

Using just “{{ forum.tagline }}” works so it seems that the problem is
limited to firstof statement. The only solution I found was doing the
firstof logic in python code and removing firstof from template.

The bottom line is that firstof is useless if you ever hope to use it
with non-ascii strings.

I’ve [opened a
bug](http://code.google.com/p/googleappengine/issues/detail?id=2303) for
this.
