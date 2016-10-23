Id: 1954
Title: Python id3 library
Tags: python
Date: 2006-04-10T20:48:58-07:00
Format: Markdown
--------------
I needed to write a program to automatically cleanup some of ID3
information in my mp3 files. I also wanted to do it in Python, since of
all languages that I know well, I'm most productive in python.

After some research I concluded that the best but little known option is
[mutagen](http://www.sacredchao.net/quodlibet/wiki/Development/Mutagen)
library.

It's written in pure pyth0n, can both read and write ID3 tags as well as
FLAC and Vorbis tags and has a nice API. It's largely undocumented but
you can figure things out by looking at its code and a sample program.

I also looked at [eyed3](http://eyed3.nicfit.net/) but it has much worse
API and
[id3reader](http://www.nedbatchelder.com/code/modules/id3reader.html),
but it can only read and has just been abandoned.

According to Google mutagen isn't a popular way to process id3
information in mp3 files, but it seems to be the best.


