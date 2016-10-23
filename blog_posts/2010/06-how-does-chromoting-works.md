Id: 224001
Title: How does chromoting works
Date: 2010-06-11T12:42:55-07:00
Deleted: yes
Format: Markdowna
--------------
### Chromoting is out

A desktop remoting solution for upcoming Chrome OS (cleverly dubbed
“chromoting”) has just been
[outed](http://www.engadget.com/2010/06/11/google-adding-chromoting-remote-desktop-functionality-to-chome/).
Is it a myth or reality? If reality, how does it work?

### How chromoting works

We don’t have to speculate. The first cut of chromoting code has landed
in public Chrome source code tree and you can [look at
it](http://src.chromium.org/svn/trunk/src/remoting/) yourself.

A look at its source code shows that remoting works exactly as you would
expect it to work:

-   on the host side, a process periodically captures the desktop as a
    bitmap at some interval (e.g.
    [capturer\_win.cc](http://src.chromium.org/svn/trunk/src/remoting/host/capturer_win.cc)
    on Windows host)
-   it calculates the difference between current bitmap and previous
    bitmap as a bunch of “dirty” rectangles. This is to minimize the
    amount of data sent to the client
-   dirty rectangles are compressed further. Currently there’s a “no-op”
    encoder and [vp8](http://en.wikipedia.org/wiki/VP8) encoder
    ([encoder\_vp8.cc](http://src.chromium.org/svn/trunk/src/remoting/base/encoder_vp8.cc))
-   this data is sent to the client (i.e. a device running Chromium OS,
    although there’s no reason this functionality couldn’t be built
    directly into Chrome browser)
-   the client decodes the data and displays it
-   the communication is done using Google’s [Jingle
    protocol](http://en.wikipedia.org/wiki/Jingle_(protocol)), which is
    based on XMPP and already used in Google’s IM and video chat clients

As far as “remote desktop” solutions go this is the simplest thing that
can possibly work and is very similar to what
[VNC](http://en.wikipedia.org/wiki/Virtual_Network_Computing) does.

On Windows more sophisticated solutions exist that hook GDI calls,
transmit the calls over the wire and reply them on the client. This is
usually faster (you need to send less data) but much more complicated to
implement, not cross-platform and possibly not as useful those days when
Windows apps increasingly use non-GDI technologies like WPF.

The source code also has a beginning of standard browser [NPP
plugin](http://src.chromium.org/svn/trunk/src/remoting/client/plugin/))
so it’s possible they’re planning a plugin for other browsers (like IE,
FireFox and Safari) and native clients for each major operating system.

### A side note on code reuse

Chromoting is a good example of code reuse. A communication layer and a
good compression for images are hard problems, so it’s good that the
project chose to re-use existing, mature technologies (Jingle and vp8).
It also shows that even if code is open-source and liberally licensed,
it’s beneficial for a corporation to own it (both jingle and vp8 code is
(now) developed and evolved mainly by Google)

### The future

The code is clearly in its early stages. Only Windows host is working
(Mac and Linux parts are just stubs). There doesn’t seem to be support
for capturing windows of a specific application (which might be a useful
functionality) etc. But I’m sure it’ll evolve quickly.
