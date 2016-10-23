Id: 2015
Title: Exporting data from EverNote
Date: 2009-02-20T01:33:37-08:00
Format: Markdown
--------------
I’ve been using [EverNote](http://evernote.com/) for some time. It’s
good piece of software, but for me it lacked some unnamed quality.

I decided to use my own website for this kind of note-taking, with
additional benefit of my notes being published on the web, hopefully for
the benefit of others as well.

There was a problem of existing EverNote notes so I wrote [this python
script](http://github.com/kjk/web-blog/blob/6a5161d2a362f2b4a35688235f9d00445f7f7f23/scripts/evernote-to-file.py)
to extract the notes and then import them as posts to this site.

It was surprisingly easy: the script is less than 150 lines of python.
EverNote team had the good sense of using sqlite for their database
storage. A little bit of reverse-engineering of the database file with
[SQLite Manager](https://addons.mozilla.org/en-US/firefox/addon/5817)
revealed that the author liked the letter Z and that the schema is
straightforward.

If you plan to use the script beware that it’s more a template from
which to start your own thing than a complete solution. It only dumps
the data I was interested in (note title, creation date, tags and
content) in a format that was simple for me to use.

EverNote also has well-documented server API, but using it seemed more
complicated than digging things out of sqlite database. And would be
many, many times slower.
