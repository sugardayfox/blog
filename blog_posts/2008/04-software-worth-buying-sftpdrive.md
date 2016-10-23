Id: 1994
Title: Software worth buying - SftpDrive and ExpanDrive
Tags: software
Date: 2008-04-15T18:47:17-07:00
Format: Markdown
--------------
[SftpDrive](http://www.sftpdrive.com/) and
[ExpanDrive](http://www.magnetk.com/expandrive) is really the same
program, except named differently for different OSes. SftpDrive is for
Windows and ExpanDrive is for Mac.

It does one thing very well: it allows mounting any ssh accessible
directory as a hard-drive/volume so that all programs running on your
Mac/PC can access it as if it was a local drive.

One use for that is as a replacement for sftp GUI client.

But what is really great is that it greatly improves interaction with
remote servers, be it an Amazon EC2 instance or a co-located server.
When I had to edit files on remote Linux servers, I used to ssh and use
Emacs for editing. Despite investing many hours in learning Emacs and
elisp (to the point that some of the keystrokes are hard-wired in my
brain) but I’m definitely disgruntled Emacs user and in my love-and-hate
relation with it I have much more hate than love. Not to mention that
for editing C/C++ code nothing comes even close to Source Insight.\
Now I can just mount a directory and transparently edit remote files
using my favorite, most familiar editor.

It’s a simple idea but the brilliance of SftpDrive is that it works
really, really well and on sufficiently fast connection I can
comfortably edit big projects. This might seem like an obvious property
of such program but my experience with Samba (either when accessing my
NAS over wireless Wi-Fi N network or accessing a filesystem exported
from Linux running on Windows under VMWare) shows that this obvious
property is lacking and trying to edit large code base in Source Insight
on a samba-exported files is an exercise in restrain from chewing your
own hand when Source Insight locks up for painful seconds when you’re in
the middle of typing code.

Two thumbs up for SftpDrive/ExpanDrive. If only they also supported
mounting S3 as a filesystem.
