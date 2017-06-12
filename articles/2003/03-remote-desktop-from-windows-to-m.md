Id: 561
Title: Remote desktop - from Windows to Mac OS X
Date: 2003-03-02T17:33:00-08:00
Format: Markdown
Deleted: yes
--------------
I just tested VNC-ing from my Windows XP box to my Mac OS X. VNC is an
open-source (GPL-ed) multi-platform remote display application.

VNC seems to be the most frequently forked application. There are:

-   the original VNC from [AT&T/University of
    Cambridge](http://www.uk.research.att.com/vnc/), no longer developed
    (I think) since AT&T shut down the research center
-   [realVNC](http://www.realvnc.com/), maintained by original authors
    of VNC
-   [tightVNC](http://www.tightvnc.com/), which focuses on improving
    compression (which improves overall quality of the software)
-   [tridiaVNC](http://www.tridiavnc.com/), Just Another Version with
    commercial PRO version (wonder how can they do it without breaking
    GPL), doesn't seem to be maintained anymore
-   [UltraVNC](http://ultravnc.sourceforge.net/), which focuses on
    Windows goodies
-   [OSXvnc](http://www.redstonesoftware.com/osxvnc/) is a server for
    Mac OS X
-   [even more links](http://ultravnc.sourceforge.net/links.html) for
    VNC-inspired software

I used OSXvnc as the server and tried realVNC, UltraVNC and tightVNC as
clients on Windows XP. I'll stick with tightVNC as it seems to be most
responsive (I'm using it with `-compresslevel 9 -quality 9` command line
settings).

So the question: **does it work?** Barely.

The quality mostly depends on the speed of your network connection. I
was VNCing over my home wireless network (UltraVNC reported the speed of
around 2400 kbit/s) and the UI redraws so slowly that working is no
pleasure at all. It works but is no fun. It has to be said that
Microsoft did much better job with their remote desktop software.
