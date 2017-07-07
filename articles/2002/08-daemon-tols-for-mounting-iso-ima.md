Id: 1383
Title: Daemon tools for mounting iso images
Tags: software
Date: 2002-08-23T03:04:59-07:00
Format: Markdown
Status: deleted
--------------
I have a new tool in my toolbox: [Daemon
Tools](http://www.daemon-tools.com/daemon_tools.htm). It allows you to
mount CD ISO images as a hard drive. Read-only, of course. What's the
use for that? My situation is: I have a RedHat 7.2 server running my web
site. I need to administer it remotely which, among other things, means
that if I need to install some additional RPM package, it's a pain
because I don't have access to installation CDs on this machine. Instead
of looking for the rpm on the net (which is little fun if all you have
is a command line) I mount the original RedHat ISO images on my main
Windows machine using Daemon Tools, I copy desired files to a remote
computer using [WinSCP](http://winscp.vse.cz/eng/) and use ssh client
[Putty](http://www.chiark.greenend.org.uk/~sgtatham/putty/) to login
remotely and install the packages.
