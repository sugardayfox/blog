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

I’ve written 4 web applications in Go, they’ve been running in
production for several months so I feel justified in publishing my
opinion.

In the past I wrote web applications in Perl, PHP, Python (web.py,
Tornado, App Engine) so those are the technologies I campare Go to.

Of the 4 websites, [AppTranslator](http://www.apptranslator.org) is a
web service for crowd-sourcing translations for software. It was written
completely from scratch.

[QuickNotes](http://quicknotes.io/) is a note-taking application. Also
written from scratch.

[Fofou](http://forums.fofou.org) is a simple forum. It is a port of an
earlier version I did in Python for App Engine. Finally, [this web
site](//blog.kowalczyk.info) is a blog engine (also a port of an
earlier App Engine version in Python).

One reason to migrate from App Engine to my own server was to save
money. At my levels of traffic
(~3 requests per second) I was paying ~$80/month, mostly for
frontend instance hours.

Another reason was to do more complex processing. App Engine is great
but being a fully managed service, it limits what you can do.

Finally, I wanted to see how Go will handle a real project.

I used to run them all on a single $60/month [Kimsufi 24](https://www.kimsufi.co.uk/) dedicated server but now I’m
deploying to [CoreOS running on $5/month Digital Ocean](/article/5/Blueprint-for-deploying-web-apps-on-CoreOS.html).

Things you need for building a web application
==============================================

What functionality is typically needed for writing a web application and
how does Go support it?

Http server
-----------

A very capable http server is part of standard library in package
([net/http](https://golang.org/pkg/net/http/)).

It parses incoming http requests and provides a way to send http
responses back.

It supports http/2 out of the box.

Url routing
-----------

High-level web frameworks provide url routing i.e. a way to say “call
this function to handle this url”. Go has a built-in simple router. I
also use [mux](http://www.gorillatoolkit.org/pkg/mux), which is built
on top of the built-in router.

Templates
---------

A lot of what web server does is returning html which is greatly
simplified by using templates. Go has a powerful
[html/template](https://golang.org/pkg/html/template/) library for that.
It's roughly equivalent of Django or Tornado templates from Python world.

There are other templating libraries for Go but I found html/template
good enough for my needs so I didn’t try the alternatives.

Cookies
-------

To generate cookies that cannot be spoofed or hi-jacked, I used
[securecookie](http://www.gorillatoolkit.org/pkg/securecookie) library

Databases
---------

Go has a [database/sql](https://golang.org/pkg/database/sql/) package
in standard library and mature drivers for all popular databases
([PostgreSQL](https://github.com/lib/pq), [MySQl](https://github.com/go-sql-driver/mysql) and [many others](https://github.com/avelino/awesome-go#database-drivers)).


Oauth
-----

There are three oauth libraries that I know of. I used [garyburd/go-oauth](https://github.com/garyburd/go-oauth) and [golang/oauth](https://github.com/golang/oauth2). They both work and I implemented
login via Twitter, GitHub and Google.

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
[XML](http://golang.org/pkg/encoding/xml/).

Modern web application are often implemented as JSON API server on
the backend and Single Page Application on the front-end that generates
HTML based on JSON returned from the server.

Returning JSON responses based on data from Go structs is very easy.

S3 access, support for zip files
--------------------------------

This one is not universally needed. My web apps do backups by uploading
data to s3, often as .zip archives.

Go has support for creating and decompressing zip and tar archives in
the standard library. For s3 support I use
[official AWS SDK](https://github.com/aws/aws-sdk-go).

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
regardless of the language used but Go makes it relatively easy.

I test my code locally on Mac but deploy to Linux-based server.

Thanks to Go's support for cross-compilation I can build a linux
binary on mac.

In the past I would copy the binary and other necesssary files to
Linux server using Fabric or Ansible scripts.

Those days I pacakge them as Docker containers and [deploy that
using simple shell scripts](/article/5/Blueprint-for-deploying-web-apps-on-CoreOS.html).

I find it crucial for productivity to quickly deploy new versions.
THat's why I write deploy scripts early in the process.

Server config and misc thoughts
===============================

In the past I would run multiple services on a single Linux server,
using nginx to multiplex.

Those days for simplicity I use a single server for each program.
That makes setup and deployment easier (no need for nginx anymore).

Show me the code
================

The source for all three projects is publicly available on Github, using
liberal BSD license: [App
Translator](https://github.com/kjk/apptranslator),
[Fofou](https://github.com/kjk/fofou),
[blog](https://github.com/kjk/blog)

Feel free to learn from the code or use it in your own projects.

Parting thoughts
================

I think writing non-trivial web services is a sweet spot for Go.

Most of the needed functionality is part of standard library and
3rdp party libraries are plantiful for more esoteric functionality.

Writing in Go is almost as fast and fluent as writing in Python but the
code is order of magnitude faster and uses less memory.
