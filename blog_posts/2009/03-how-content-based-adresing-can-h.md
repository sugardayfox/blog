Id: 8074
Title: How content-based addressing can help web performance
Date: 2009-03-05T16:54:31-08:00
Format: Markdown
--------------
### The problem with HTTP caching

Currently an HTTP resource (e.g. an image, a JavaScript file, a CSS
script) is uniquely identified by a uri. Unfortunately a web browser
doesn’t know if a given resource has changed, so in the simplest case,
it has to re-download it over and over again. This is very inefficient,
so we came up with ways to tell the browser to avoid re-downloading
content.

One method is by setting Expires HTTP header by which a web server can
tell the browser for how long can it cache a given resource. The problem
with this approach is that if the resource is changed before the
expiration date, the browser will see an outdated version.

Another method is [ETag](http://en.wikipedia.org/wiki/HTTP_ETag). When
returning a resource web server also returns an ETag header which
uniquely identifies the resource. If subsequent response contains the
same etag, the browser might skip downloading the content. This still
requires making a small request.

None of this methods works globally i.e. if two websites use the same
version of jQuery library, even if they both instruct the browser to
cache them, the browser still ends up downloading two versions. Google’s
[AJAX Libraries API](http://code.google.com/apis/ajaxlibs/) is an
attempt to solve this problem for the most popular web libraries.

What if there was a backwards-compatible way, efficient way of caching
popular resources that worked globally (across multiple websites)? In
theory it’s possible and not that hard.

### Enter content-addressable resources

The idea of content-addressable data has been recently used to great
benefit in e.g. [git](http://git-scm.com/) or [BitTorrent
protocol](http://bittorrent.org/beps/bep_0003.html).

The idea is simple: we can uniquely identify any piece of data by
calculating sha1 (or some other cryptographically secure hash) of its
content.

All we need is a way to tell the browser the sha1 hash of a resource. If
a resource with this hash is already cached locally, the browser doesn’t
need to re-download it.

Compared to existing solutions, it has following advantages:

-   it works globally. If 10 websites host the same version of jQuery
    library, the browser only need to download it once
-   it doesn’t have the ‘oops, we set the expiration date too far into
    future’ problem that setting Expires header
-   it’s more efficient than ETag because it doesn’t require making any
    requests

### How to tell the browser about the hash

The simplest way would be to define another attribute e.g. sha1hash that
could be added tags that refer to resources, like script, a, link etc.

I’m not a web technologist and I don’t know this particular solution
would be acceptable, but I’m sure there is a way.
