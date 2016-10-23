Id: 759002
Title: Showing html from memory in embedded web control on windows
Tags: note,win32,programming
Date: 2011-11-30T04:30:12-08:00
Format: Markdown
--------------
This came up while working on
[SumatraPDF](http://blog.kowalczyk.info/software/sumatrapdf/) and took
me a while to figure out, so I’m documenting it for the posterity.

The problem: an easy way to show HTML in a windows program is by hosting
(embedding) a web browser COM control in win32 HWND. To display HTML
that lives on some http server or is a file on disk is easy: call
IWebBrowser2::Navigate2().

What if the HTML is just a piece of data in memory? Easy, call
IWebBrowser2::write(). But what happens if that HTML refers to other
things e.g. a CSS file or a PNG image? Well, that becomes really
complicated. Solution I’ve settled on can be summarized as:

1\. register custom IInternetProtocol/IInternetProtocolInfo/ via custom
IClassFactory given to IInternetSession::RegisterNameSpace(). For
reasons that seem like a bug to me, it has to be a protocol already
known to IE (I’ve chosen “its”) even though it would be much better if
it was my own, unique namespace.

2\. feed html data via custom IMoniker through IPersistentMoniker::Load()
and make sure that IMoniker::GetDisplayName() (which is a base url
according to which relative links in provided html will be resolved)
starts with that protocol scheme (in my case “its://”). That way
relative link “foo.png” in the html data will be its://foo.png to IE
which will make urlmon call IInternetProtocol::Start() and
IInternetProtocol::Read() to ask for the data for that url.

This is all rather complicated, you can look at the actual
(BSD-licensed) code:
[HtmlWindow.cpp](http://code.google.com/p/sumatrapdf/source/browse/trunk/src/utils/HtmlWindow.cpp)

And now a few related rants.

### Bad design, bad documentation from Microsoft.

Hosting web browser control is a maze of interdependent, interconnected
COM interfaces that really like to interact via seven layers of
indirection. There are [5 ways to load HTML from
memory](http://qualapps.blogspot.com/2008/10/how-to-load-mshtml-with-data.html),
none of which allowed me to easily resolve linked resources. That’s just
bad design.

Microsoft wrote a bunch of documents documenting the basic usage of the
embedded browser control but the documentation is far from satisfactory.
They recognized the need to load HTML from memory (it’s well documented)
but I just can’t believe that they didn’t consider the scenario of
linked documents (HTML is, after all, a format based on links). Getting
links to embedded resources is not easy (unless I missed the easy way to
do it) and not documented at all.

### StackOverflow not so useful for involved question.

After battling the issue for many hours, I did what a modern programmer
is supposed to do in such case: [I asked a question on
StackOverflow](http://stackoverflow.com/questions/8265887/how-to-provide-image-data-for-embedded-web-control-in-c).

That turned out to be a big disappointment. I got two responses, neither
of which presented acceptable solution. In fact, out of [5
questions](http://stackoverflow.com/users/2898/krzysztof-kowalczyk?tab=questions)
that I asked on StackOverflow, only one got what I consider to be
complete, expert answer.

I do tend refer to StackOverflow only after reaching the state of
desperation and exhausting all my other options (i.e. web and code
searches) so my questions are probably a bit harder than an average
question but still 20% answer rate is bad.

StackOverflow can give you an answer to an easy question literally in
minutes, but the hard ones are often not answered well at all.

### How will I live when Google’s code search is gone?

The only reason I was able to figure this solution was by using Google’s
code search to find out how other programs do it but Google just
announced that they’re [shutting it
down](http://googleblog.blogspot.com/2011/10/fall-sweep.html)

Not that I think it was a great service but it’s certainly better than
anything else in that domain. What will I use the next time I have a
hard time understanding someone’s clunky APIs?
