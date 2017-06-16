Id: 686002
Title: Introducing Volante - a database for C# (.NET)
Tags: software,volante,.net,c#
Date: 2011-09-22T17:12:27-07:00
Format: Markdown
--------------
[Volante](/software/volante/database.html) is
a small, fast, object-oriented, embeddable database designed for
seamless integration with C\# (and other .NET languages). Today I’ve
made a first public release.

Almost every program needs to persist some data. There are many ways to
do that: serialize data as XML or JSON, use SQLite etc.

I’ve been recently writing desktop .NET applications in C\# and existing
(free) options didn’t meet my needs well.

The closest best solution is SQLite, but it doesn’t integrate with C\#
well: you have to convert your object to/from SQL tables.

If I was the world’s toughest programmer I would write an
object-oriented database engine from scratch and it would be designed
from the beginning for seamless integration with .NET framework.

Thankfully, I didn’t have to do that. What I needed already existed in
the form of [Perst](http://en.wikipedia.org/wiki/Perst) project.

There was one wrinkle: while early versions of Perst were under BSD
license, with version 2.50 the code was acquired by a company McObject
and is now distributed under GPL and those who can’t use GPL can
purchase commercial license from McObject.

Not so great for one person who doesn’t yet make money from his
software.

I decided to adopt the Perst code base. I picked the latest 2.49 version
that was still licensed under BSD (copyright cannot be change
retroactively) and I’ve spent the last couple of months writing
comprehensive documentation, writing tests, fixing bugs discovered by
tests, modernizing the code base.

Today I’ve reached a point where I’m comfortable releasing this code
publicly as version 0.9.

I’ve retained the BSD license of early Perst versions so the code is
free to use in both open-source and commercial projects.

Volante database serves the same niche as SQLite: an embedded database
engine for your desktop C\# applications. Like with SQLite, the database
is in a single file.

There are significant differences from SQLite.

.NET is an object-oriented framework. Volante is an object-oriented
database to offer the best integration with .NET. Volante uses B-Trees
to implement indexes, which allows quickly finding objects with desired
properties.

Volante is extremely small: Volante.dll is only 180 KB.

I distribute Volante.dll for Microsoft’s .NET (works in .NET 2.0 and
later) but it can also be compiled from sources and used under Mono.

I was dogfooding Volante from day one in my three .NET applications and
it’s been performing great.

I hope you’ll find it useful too.
