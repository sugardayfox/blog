Id: 424001
Title: XML is really, really slow
Tags: programming
Date: 2011-03-10T18:25:56-08:00
Format: Markdown
--------------
It’s about 8 years late but W3C finally produced a standard for
[efficient XML
representation](http://www.readwriteweb.com/archives/new_xml_standard_for_super-fast_lightweight_applic.php).
They call it somewhat opaquely [Efficient XML
Interchange](http://www.w3.org/TR/2011/REC-exi-20110310/) (EXI). I guess
it sounds better than what we used to call such things back in the day:
binary XML (which is a better description what it’s about).

The money quote from the announcement:

> They’ve achieved over 100-fold performance improvements…

One way to interpret this is: wow, those guys are really smart.

There is much simpler explanation: XML is really, really slow.

Speed comes from architecture
-----------------------------

I’m somewhat performance oriented in my programming work. One of the
reasons I disliked the popularity of XML was that I saw how often it
blinded people to engineering realities. Choosing XML was often a reason
for dramatic performance issues that then had to be heroically
recovered.

Those performance problems were, however, mostly self inflicted. XML is
only one of the possible ways to store or exchange data but it certainly
is one of the slowest.

You cannot get 100x speed up over technology that wasn’t incredibly
inefficient to begin with.

Sometimes you don’t have the choice (when you have to inter-operate with
systems that only offer XML) but I’ve seen many cases where people did
have the choice and made the wrong one.

XML isn’t such a hot buzzword anymore so I don’t see that problem as
often but I still do see it.

For example, more than one text editor decided to store the syntax
highlighter definitions as XML while it can be trivially stored in
simple text format that can be parsed much faster, using less memory
and, as a bonus, can actually be edited by human beings.

A bigger point is: speed comes from architecture. It’s not that choosing
the right architecture will make you dramatically better but choosing
the wrong architecture will make you dramatically worse.

Choosing XML over more efficient ways of storing data is just one, but
particularly frequent, example of that rule.

Binary XML (EXI) is too little too late
---------------------------------------

As to the standard itself: don’t pay attention to it. It’s too little,
way too late.

About 8 years ago I did my small part in expanding XML capabilities of
SQL Server. At the time XML was hot so Oracle, Microsoft, IBM raced to
add native XML handling to their databases. It was going to be the next
big thing. Possibly even bigger than big.

Trust me, the ridiculous inefficiency of XML wasn’t lost on developers
working on XML technologies 8 years ago. Coming up with a more efficient
binary XML format is easy. Microsoft had its version (if not several of
them) and other companies had theirs.

The real problem was: politics prevented people agreeing on a standard
so no standard emerged.

There was an eruption of creating standards based on XML (SOAP, XML
Schema, XQuery). If you don’t know what those terms mean it’s because
they all failed (despite the fact that everyone was convinced they’re
going to be the next big thing).

In hindsight it was a terrible mistake to work on those big standards of
speculative value but not solve a real problem people already had: an
efficient, binary, standard format for storing XML.

If W3C came up with EXI 8 years ago, maybe people wouldn’t feel the need
to invent Protocol Buffers or Thrift and it would have won.

But solving this problem today is almost comically late. It has no
chance of adoption.

If you’re planning to use a custom, binary way of storing XML, using EXI
is probably better but there’s no way EXI will become so universally
supported as XML and universal support (in software libraries, books
etc.) is really the main thing XML has going for it (speed or
human-readable syntax, on the other hand, are not XML’s strengths).
