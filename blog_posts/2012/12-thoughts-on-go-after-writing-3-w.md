Id: 1441010
Title: Thoughts on Go after writing 3 websites
Tags: go,programming
Date: 2012-12-31T15:30:45-08:00
Format: Markdown
--------------
Go for writing web servers
==========================

**Summary**: in my experience Go is a good language for building
websites/web servers.

It’s easy to get excited about new technology like Go. The question is:
how does it stand up to scrutiny after daily use?

I’ve written 3 web applications in Go, they’ve been running in
production for over a month so I feel justified in publishing my
opinion.

In the past I wrote web applications in Perl, PHP, Python (web.py,
Tornado, App Engine) so those are the technologies I campare Go to.

Of the 3 websites, [AppTranslator](http://www.apptranslator.org), is web
service for crowd-sourcing translation for software and was written
completely from scratch.

[Fofou](http://forums.fofou.org) is a simple forum and is a port of an
earlier version I did for App Engine. Finally, [this web
site](http://blog.kowalczyk.info) is a blog engine (also a port of an
earlier App Engine version).

One reason to migrate from App Engine to my own server was to save
money. At my levels of traffic
(~3 requests per second) I was paying ~$80/month, mostly for
the frontend instance hours.

Another reason was to do more complex processing (App Engine is great as
long as you don’t have to do something that App Engine doesn’t support).

Finally, I wanted to see how Go will handle a real life project. The
best way to test a new technology is on a project with a predictable
(and relatively small) scope.

I used to run them all on a single $60/month [Kimsufi 24](http://www.kimsufi.co.uk/) dedicated server but now I’m using
$10/month [DigitalOcean](https://www.digitalocean.com/) VPS. I’m using
latest Ubuntu for the OS.

Things you need for building a web application
==============================================

What functionality is typically needed for writing a web application and
how does Go support it?

Http server
-----------

A very capable http server is part of standard library
([net/http](http://golang.org/pkg/net/http/)).

It parses incoming http requests, provides an easy way to send http
responses back.

Url routing
-----------

High-level web frameworks provide url routing i.e. a way to say “call
this function to handle this url”. Go has a simple built-in router. I
also use [mux](http://www.gorillatoolkit.org/pkg/mux) which is built on
top of the built-in router.

Templates
---------

A lot of what web server does is returning html. Constructing this html
is greatly simplified by using templates. Go has a powerful
[html/template](http://golang.org/pkg/html/template/) library for that.
To me it seems roughly equivalent of Django or Tornado templates.

There are other templating libraries for Go, but I found the above
built-in package satisfactory, so I didn’t even try to use them.

Cookies
-------

Basic support for cookies is part of built-in http library. To generate
cookies that cannot be spoofed or hi-jacked, I used
[securecookie](http://www.gorillatoolkit.org/pkg/securecookie) library

Databases
---------

There are Go libraries for all of the popular databases (MySQL,
PostgreSQL, MongoDB, Redis etc.) but I haven’t used them.

I used what I call NoDB approach i.e. I wrote a very simple storage
system that uses text files. Data is kept in memory and persisted in an
append-only file.

This wouldn’t be the right approach for services that require more
sophisticated functionality but was good enough for my needs and didn’t
take long to implement.

Oauth
-----

There are three oauth libraries that I know of. I used [this
one](https://github.com/garyburd/go-oauth). I didn’t have any particular
reason to choose this one over the others. I only needed it for
implementing Twitter-based authentication, this library worked so I used
it.

Generating atom (rss) feeds
---------------------------

One cannot respect a blog engine that doesn’t provide full-text rss
feed. I couldn’t find an existing package so I build a simple (and
small) library for [generating atom
feeds](https://github.com/kjk/atomgenerator)

JSON and XML support
--------------------

Go has a built-in support for
[JSON](http://golang.org/pkg/encoding/json/) and
[XML](http://golang.org/pkg/encoding/xml/).\
There are APIs for raw parsing/serialization. If you know the shape of
the data, you can serialize (marshal, in Go’s parlance) data structures
to JSON or XML and de-serialize (unmarshal) from JSON or XML into a
struct.

S3 access, support for zip files
--------------------------------

This one is not universally needed. My web apps have built-in backup
functionality which stores data, sometimes in the form of a .zip file,
in s3.

Go has support for creating and decompressing zip and tar archives in
the standard library. For s3 support I use
[goamz](http://github.com/crowdmob/goamz).

Unit tests
----------

This is not specific to writing web server - all your important code
should have unit tests.

Go has a [built-in API and tool
support](http://golang.org/doc/code.html#Testing) for writing and
running tests.

Deployment
==========

Deploying a new version of your code to the server is a pain in the ass
regardless of the language used.

I wrote relatively short deployment script using
[Fabric](http://fabfile.org) (which is a python library and a tool for
running deployment scripts). It copies the source files to the server,
compiles them on the server, runs unit tests, shuts-down existing
instance of the service and launches the new version. It stops if
there’s anything wrong along the way.

It’s really important to be able to quickly deploy new versions of the
software so those days I would write the deployment script as the first
thing in the project. Doing all those deployment steps manually would be
very annoying.

Server config and misc thoughts
===============================

The overall setup of the server is pretty standard: each service is
running as a separate process on a local port. Nginx is running on port
80 and proxies the traffic to a given service based on Host header.

Show me the code
================

The source for all three projects is publicly available on Github, using
liberal BSD license: [App
Translator](https://github.com/kjk/apptranslator),
[Fofou](https://github.com/kjk/fofou),
[blog](https://github.com/kjk/web-blog)

Feel free to learn from the code or use it in your own projects.

Parting thoughts
================

I think writing non-trivial web services is a sweet spot for Go.

Most of the needed functionality is part of standard library. For almost
everything else there are 3rd party libraries.

Writing in Go is almost as fast and fluent as writing in Python but the
code is order of magnitude faster and uses less memory.
