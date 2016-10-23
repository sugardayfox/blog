Id: 9002
Title: Ideas for software
Tags: idea
Date: 2009-02-21T15:07:35-08:00
Format: Markdown
Deleted: yes
--------------
### Port SelfControl to Windows

[SelfControl](http://visitsteve.com/made/selfcontrol/). Would have to
hook up into network stack and block sites, showing a page that is e.g.
like [this](http://phylab.mtu.edu/~nckelley/Focus/)

### Port ShoveBox to windows

[ShoveBox](http://www.wonderwarp.com/shovebox/)

### Posterion

Posterous desktop client for Windows.

### Tumbletor

Tumblr desktop client for Windows.

### Internet Traffic Spy

Snoops tcp traffic and remembers how much data was sent to/received from
a given ip/site and shows stats on that (daily, weekly, monthly).

### DNSSpy

Snoops tcp traffic, measures time taken to resolve DNS queries and plots
it over time. Also remembers which domains were resolved and can show
it.

### Redis UI

Can show graph of stats and browsing of keys and their values.

### korev

Code review tool for mac/windows. Helps in two scenarios:

-   previewing local changes before committing them to a repository.
    Allows entering commit message and committing directly from the app
-   requesting a code review. Uploads changes to a website, sends
    request e-mails etc.

Support subversion, git, mercurial.

### Upload to various image-sharing services

The easiest way (on Windows) to upload images to various image-sharing
services like [yfrog.com](http://yfrog.com/), twitpic, flickr, imagur,
zoomr, s3, smugmug etc. (google for “photo sharing”, “image sharing”)

### Port skitch to Windows

It could either be done as adding skitch functionality to the image
uploader app or tied to a website built just for that (which duplicates
skitch.com functionality)

### Pixel art editor in a browser

Written in Silverlight.

### Language learning website

[More info](/article/Idea-for-language-learning-website.html)

### HTTP inspection/debugging

Via pocket capture. Save all requests/responses, allow inspecting
headers etc. Automatically diagnose problems in the html, css,
javascript etc. (like redbot)

### Port soulver to WPF/Android

[Soulver](http://www.acqualia.com/soulver/)

### Clone of Tweetie in WPF

[Tweetie](http://www.atebits.com/tweetie-mac/)

### Simplenote desktop app for Windows/Android

[Simplenote](http://simplenoteapp.com/extras). Could be a port of
[Nottingham](http://clickontyler.com/nottingham/)

### A web based PDF reader

Use Sumatra as a marketing tool. Allow uploading files directly from
Sumatra or from a separate application (for bulk uploads).

It would priced like DropBox, in a way that makes profit on top of s3
prices. It would be free for a small number of documents (15-25) and
then \$9.99/month or \$19.9/month depending on number of files (number
of files TBD).

Needs iPad-optimized view.

### Software version of Printable CEO Task Tracker

[Printable CEO Task
Tracker](http://davidseah.com/blog/the-printable-ceo-part-ii-much-to-do-about-task-tracking/)

### Notenik

For mac, easy note taking with web based synchronization. The idea is to
focus on the client-side editor. Could be extended into a blog editor.

### VisualAck

For mac, finish it.

### Clipboard on-line backup

For mac/windows. Monitors clipboard and saves everything to a website,
making it easy to copy&paste between computers.

### [UrbanDictionary](http://www.urbandictionary.com/) iPhone client

Partner with UrbanDictionary for content and promotion, split revenues
in half.

### web-based **.cbr/**.cbz comic viewer

People would upload their **.cbr/**.cbz comic files. Using flash we
could do full screen browsing.

### Awesome on-line repository history viewer.

Main problem to solve: browsing history of checkins, figuring out who
changed that line of code when.

Secondary problem: browsing/searching the code.

Google code is pretty good for browsing, but it’s not good at searching
history of checkins (i.e. find in which checkin a given line/text was
changed).

Simplest implementation: full-text index every diff (probably with lexer
customized for code) and search.

### Search oriented BitTorrent client.

A prototype in C\# based on that BitTorrent client library. Search is at
the center of the app (like google) and plugs into one (or more) search
backends with a standardize interface. Important: search results should
be streamed i.e. displayed as they’re being returned (the way one of the
web-based searches do) so that the user doesn’t get impatient. Prototype
search backend as a proxy to a website.

### Port qemacs design to C\#

Could end up as a text viewer engine in many apps. Or as an editor
engine for some apps. Or as on-line editor.

### C\# Metakit (http://www.equi4.com/metakit/) port

Just to see if it fits well with C\#.

### Resurrect Perst

Resurrect
(http://web.archive.org/web/20041011051317/http://www.garret.ru/\~knizhnik/perst.html)
starting from the latest Java source I can find (so far 1.5), port to
C\#, host on gihub.

### Website that allows uploading PDFs and viewing in the browser

A simple python app (using webpy.org and sqlite that is part of
python2.5 for database?) that allows uploading of PDF files. If a file
is determined to be a valid PDF, it’s saved to s3 under its sha1 name +
.pdf suffix in some of my buckets (kjkpub under /pdfs directory?). We
remember basic info about submitted PDF (original filename, submition
time, sha1 name) in database. We then spin a process that renders each
page as a \*.png file (possibly making sure that grayscale images are
saved as 8-bit to save space) and save each .png under their sha1 name
on s3. We remember which pages belong to which PDF and which version of
the program rendered them. We allow viewing of PDFs in the browser by
showing .png files.

This is meant for submitting broken PDFs, so on submission page there’s
also a text box that allows entering description of whats wrong with the
file.

### Access computer files remotely.

From iPhone or through web. Has already been done, but not necessarily
right (orb, avvenue).

### Dictionary for iPhone.

WordNet 3.0 + Wikitionary-based dictionaries. All dictionaries for one
price of \$9.99. Dictionaries can be downloaded from the web.

### Torrent inspector

Dumps info about a torrent file. In silverlight or on App Engine.

### On-line chm viewer

People can upload \*.chm books and view them on-line. Remembers when a
given page was read (viewed) by a given user.

### Website for US comic lovers.

Base of functionality would come from cover scans extracted from DCP
releases automatically downloaded from bittorrent network. It would link
to places where the comic could be bought. It would try to spider misc
metadata (authors, title, serie number, release date (month/year),
publishing house). It would have pages for authors, series. It would
allow editing metadata by people in a wiki way. Accounts for people.
Leaving comments. Standard social features. Show most active
contributors to the site (those who edit metadata etc.). Links to amazon
and what not. Ability to upload new scans. Links to reviews on other
websites. Estimated time: 1 month (ha, ha)

### xobni for gmail and/or mac

http://www.techcrunch.com/2008/01/09/xobni-the-super-plugin-for-
outlook/

### “Time Yourself” website.

Allows to enter a task, start/stop timer on that task + shows history of
tasks. Market as productivity/anti-procrastination tool. Write about
time boxing technique. Other names/slogans: “What I did”, “What I do”,
“Timeboxer”, “Timing the Monster”.

### del.icio.us windows native client.

Needs to be fast, small and good looking.\
Resources:

-   [.net library docs](http://netlicious.sourceforge.net/) and
    [code](http://sourceforge.net/projects/netlicious/develop)
-   [netlicious desktop client](http://www.procanta.com/), for windows,
    pretty awful
-   [pukka](http://codesorcery.net/pukka), simple mac client (\$17)
-   [delicious api docs](http://delicious.com/help/api)
-   [decoding/encoding urls without web client
    dll](http://www.west-wind.com/weblog/posts/617930.aspx) (which is
    not part of client profile)
-   [another post about
    that](http://blog.coditate.com/2008/12/version-of-systemwebhttputility-for-net.html)

### longhand for Windows

Build [longhand](http://scottfr.googlepages.com) for Windows (or as a
webpage in Silverlight)

### TextMate clone for Windows

(already exists in e and that sublime text)

### super-file-manager for Windows

With support for scp and s3. Basic functionality would be free, scp and
s3 would be in PRO version.

### Supermemo for iphone/android

Flash card program for web with iphone/android/blackberry client for
repetitions

### Subversion/git client

For Windows and/or Mac

### Desktop on-line storage aggregator

Windows/mac app that would aggregate free on-line storage accounts.
Possibly with a driver for mapping it as a local drive. Aggregates
dropbox, microsoft’s storage, google’s storage etc.

### super-console

For Windows or mac. Some ideas from
[DTerm](http://www.decimus.net/dterm.php). Also ssh from putty.

### Lighthouse desktop client

Like [Lighthouse
Keeper](http://www.mcubedsw.com/software/lighthousekeeper) but for
Windows.
