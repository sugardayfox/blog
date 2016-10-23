Id: 2017
Title: BitTorrent-based, large file distribution for HTTP
Tags: idea,bittorrent
Date: 2010-04-18T15:35:18-07:00
Format: Markdown
--------------
Serving large files to many users at the same time is not a solved
problem on the web. Unless, of course, you have enough money to pay
Akamai to host your files. But even Akamai [isn’t
infallible](http://www.techcrunch.com/2009/01/21/the-day-live-web-video-streaming-failed-us/)

The basic problem is that no matter how big your pipe to the internet
is, if enough people want to use it at the same time, you’ll eventually
run out of outgoing bandwidth. Additionally, big pipes cost big money so
only few can afford them.

The problem in general has been solved by BitTorrent protocol with a
simple, yet brilliant idea: convert downloaders into uploaders and
offload some of the outgoing bandwidth needed by the host to peers
(other downloaders). The beauty of this solution is that it scales with
the traffic: the more peers try to download your stuff, the more peers
there are to also upload your stuff. BitTorrent is wildly popular: it
has been reported that BitTorrent traffic is responsible for 150% of
total internet traffic and generates 17% of technology-related lawsuits.

The problem has not been solved yet specifically for the web. All the
pieces are almost there, they just need to be put together in just the
right order. So here’s what I think is needed to bring BitTorrent
advantages for large file distribution to the web:

-   browser integration
-   no tracker, DHT-only
-   support for web seeds
-   open-source implementation with no strings attached
-   additional bonus: streaming support

Let’s take them one by one.

### If it ain’t in the browser, it doesn’t exist

Browsers already natively support HTTP and FTP downloads (at the
minimum). It would be imperative to have support for this kinds of
downloads natively in the browser, as part of regular download manager.
Users should not even be aware they are using BitTorrent protocol.

It’s a chicken and egg problem: no files available for this kind of a
download means no incentive for browser vendors to add support for them.
But miracles do happen and Mozilla in particular has recently shown
strong leadership in moving all things web forward. Having average Joe
be able to distribute large files efficiently and cheaply, even if they
become wildly popular, would be good for the web. Humankind, even.

### Trackers? We don’t need no stinking trackers

Originally BitTorrent protocol required additional piece of
infrastructure: a tracker. Tracker is a broker - it tells downloaders
about other downloaders. This creates a single point of failure - if
tracker doesn’t work, nothing works.

A later addition to the protocol was support for DHT (Distributed Hash
Table). It does the same thing as tracker except it uses ad-hoc network
created by downloaders themselves to exchange information about who’s
downloading what.

There’s still a bootstrap problem: you have to make first DHT query
somewhere so there would still be need for a bootstrap DHT server.
Bootstrap server could be run by one or two trusted entities (e.g.
Mozilla or Google).

### Things that grow on the web

In BitTorrent there’s no permanent, canonical source of your data. The
original uploader becomes the first source of data and the more people
download it (and remain uploaders), the more it becomes replicated.

However, if all current uploaders (called seeders in BitTorrent lingo)
go offline, the file is gone.

On the web we do have the file available via HTTP. Support for web seeds
simply means that location of this file is known (by being encoded
inside .torrent file) and software is able to use it even if there are
no other people to download from. In the worst case the performance is
no worse than plain HTTP download.

### Open web requires open-source

Open-source implementation is a must for inclusion with Firefox or
Chrome and would go a long way in convincing other browser vendors to
include it as well.

The license should be non-restrictive (BSD, MPL, Apache but not GPL or
LGPL).

Currently the most viable choices are
[libtorrent](http://www.rasterbar.com/products/libtorrent/index.html)
and [unworkable](http://p2presearch.com/unworkable/).
[libtransmission](http://www.transmissionbt.com) is not on this list
because of the weird “mostly MIT but with selected files under GPL”
license.

### How it would work - the big picture

Those who wish to enable BitTorrent support for their files would have
to create \*.torrent files for each download (possibly with a different
extension, e.g. .webtorrent, so as to not confuse those files with
regular .torrent files).

Those files could either be exposed directly and browsers supporting
them would be able to download them.

Alternatively, there could be a naming scheme e.g. that for
`/foo/bar.iso`, corresponding torrent file would be
`/foo/bar.iso.webtorrent` and browsers supporting this scheme would also
try to hit this url.

Doing additional requests, potentially for nothing, might be too much so
we could come up with some metadata to embed in html itself that would
instruct the browser that file foo.iso has corresponding
foo.iso.webtorrent.

### Bonus - streaming support

Often the large files people want to download are video files and they
aren’t actually interested in downloading them, just watching once. This
solution could be supported to support streaming as well - there already
are modifications to BitTorrent algorithm to work in streaming mode i.e.
try to deliver the file sequentially from the beginning at a desired bit
rate.

### Competing existing solution

This idea came to me while reading about
[zsync](http://zsync.moria.org.uk/) which has similar requirements for
generating additional metadata files on the server. Zsync optimizes for
subsequent updates of the same file, this idea optimizes for cheap and
fast delivery of large and popular files.

[Some websites](http://torrent.fedoraproject.org/) have special section
with torrent version of their files. This method requires more work on
the provider part and requires users to use a special client for
download. BitTorrent clients, while extremely popular, are not as
popular as web browsers and are more difficult to use.

Amazon’s s3 hosting has an option for automatic generation of .torrent
files for any file hosted with s3. It removes the need for running a
tracker and manually generating .torrent files but still requires users
to use a special client.
