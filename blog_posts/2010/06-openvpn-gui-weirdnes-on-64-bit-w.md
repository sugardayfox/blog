Id: 200001
Title: OpenVPN gui weirdness on 64-bit Windows 7
Tags: software
Date: 2010-06-06T14:42:05-07:00
Format: Markdown
--------------
This is just to document a weirdness in [OpenVPN
2.1.1](http://openvpn.net/index.php/open-source/downloads.html) on
64-bit Windows 7 so that it helps others who also encounter this
problem.

The weirdness: even though my config.ovpn file has hard-coded paths to
my certificate and keys (e.g.
`c:\Program Files (x86)\OpenVPN\config\ca.crt`) when run as
administrator, OpenVPN gui insists on looking for them in
`c:\Program Files\OpenVPN\config` directory. No amount of editing to
that file changes that.

It doesn’t happen if OpenVPN gui isn’t started as administrator, but
then it doesn’t work because it isn’t allowed to change routing tables.

The only work-around I found is to create
`c:\Program Files\OpenVPN\config` directory and duplicate the files
there. Very frustrating, but it’s better than Mac, where I couldn’t get
Viscosity to work with the same config (it connects and then immediately
drops the connection).
