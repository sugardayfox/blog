Id: 3085
Title: Music backup service idea
Tags: idea,note
Date: 2009-02-18T01:03:18-08:00
Format: Markdown
Deleted: true
--------------
Unfortunate event
=================

Has already been done by [bluetunes](http://www.bluetunes.net/),
[AudioBox](http://audiobox.fm/) and probably many others.

On-line music backup application
================================

A free backup of music files (mp3, wav (?), ogg (?) + cover images).
Make money by offering premium service (at least \$4-\$5/month) e.g.
streaming of the files through a web browser and/or native iPhone app.

A trick: I assume that non-drm mp3 files have exactly the same content
so that I only need one copy of the file for all people that backup this
file on my server =\> huge storage savings.

This might not be true if iTunes fucks with metadata inside mp3 file.

Components
==========

1\. A mac/windows desktop app

-   configures which folders to backup (and helps by scanning the whole
    hard-drive). The UI for selection of folders to backup like in
    Picasa
-   builds a catalog of files (file path + sha1 hash of the file) and
    sends that to the server in one big chunk (initial setup)
-   sends all the files to the server (if not already there)
-   detects when one of the tracked files gets deleted and offers to
    restore or archive (i.e. no longer track this file - it’ll still be
    available on the server). A separate section in UI for archived
    files.
-   detects new files in tracked folders and notifies the server about
    them & sends the files

2\. A server side of the uploader

An always-running process that:

-   can answer the question: is this sha1 already present on the server
-   accepts upload of a given file (sha1, size followed by the content
    of the file)
-   probably written in C\# or C (efficient with low resource usage)

3\. Web app

Account creation, browsing your files through the web, downloading the
files etc.

Misc ideas.
===========

Ability to share all (or some) files with people you know well. The
exact details tbd.

Sharing option 1: both parties have to have account. Both pull & push
way to initiate sharing. Push: the owner of files knows the identity (on
our server i.e. user id) of the person he wants to share his files with.
He selects which files and user id. The other person gets notified.
Pull: the person who knows the identity of the sharer requests access to
his music collection. The (potential) sharer gets notified and can
accept or reject request. After accepting, he selects which files he
wants to share (all or a subset). Push is needed to support a
streamlined process where user of the service can send an e-mail to a
person that doesn’t have an account which will explain how to get access
to the music (1. create an account 2. request sharing from the person
that initiated the sharing)

Ability to create a mix tape i.e. select N songs (N no bigger than, say,
10) with interface similar to muxtape, available under some random url.
Mix tape will expire after some time (e.g. after a week). This is so
that a user can recommend few songs to someone that they don’t know well
enough to share the songs.

iPhone app that allows streaming the whole music collection.

Cost analysis i.e. what it costs to run the service
===================================================

Major costs:

-   servers for web serving and upload management. We could start with
    just one small ec2 instance (\$72/month) but would be safer to do 2
    (one for web server, second for upload management) i.e.
    \~\$150/month
-   storing music files. I assume the required storage capacity will be
    in terabytes. Storing 1TB on S3 is \$150/month. Assuming my theory
    about (lack of) uniqueness of mp3 files is right, storage cost will
    be pretty constant, not growing very fast (except initially) and the
    cost per-customer will drop the more customers we have
-   upload costs (when doing the initial backing up). Even cheaper:
    uploading 1TB is \$100 and it’s a one-time cost
-   streaming - this will vary depending on usage. Downloading 1TB is
    \$170 on s3. Downloading 23 GB is \$3.91 so if our monthly price is
    \$4, the our breaking point would be at 23GB bandwidth used per
    user. At 192 kbits/s, 3 minutes of music is 4.2 MB. At 129 kBits/s
    24GB is, I believe, 291 hours = 36\*8hrs, so we should be safe.
    Also, going with cheaper hosting provider for bandwidth would be an
    option if that becomes a factor.
-   restoring costs i.e. downloads due to restoring the backup.
    Hopefully that won’t be a frequent event

Marketing piches
================

Your music. Safe.

The cheapest way to backup your music. (well, because it’s free).

Similar services
================

http://orb.com - also allows to share music (among other stuff)\
Pandora - 17 million users, 1.5 million on iPhone, according to
http://www.alleyinsider.com/2008/10/are-half-of-pandora-s-listeners-on-the-iphone-no
