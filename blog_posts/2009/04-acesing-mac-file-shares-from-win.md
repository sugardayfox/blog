Id: 20003
Title: Accessing Mac file shares from Windows 7
Tags: mac,windows
Date: 2009-04-11T15:40:54-07:00
Format: Markdown
--------------
Windows 7 beta 7000 can’t access files shared from Mac through built-in
samba file sharing out of the box. It’s rather infuriating, because XP
has no such problem and the error message given by windows when you try
to navigate to a folder shared from mac is cryptic and not helpful at
all. Turns out the culprit is purely how Windows 7 security is
configured by default. Good news is it can be fixed via configuration
tweaks as [described
here](http://www.pronetworks.org/forums/windows-7-and-mac-os-x-sharing-t104898.html).

In case this post goes away, here’s what it says:

Windows 7 will not work with Mac OS X Windows file sharing support by
default. If you attempt to access a folder shared from Mac OS X, Vista
will display a logon error repeatedly. The problem is that Vista, by
default, will only use NTLMv2 for authentication, which is not supported
by Mac OS X’s Windows Sharing service. The other problem is the Minimum
Session Security for NTVLM SSP based Clients.\
To get around this:

1\. In Vista, open the Control Panel\
2. Switch to “Classic” view\
3. Double-click Administration Tools\
4. Double-click Local Security Policy\
5. Or Secpol.msc\
6. Expand “Local Policies” and select “Security Options”\
7. Alternate : Type secpol.msc to get editor up then\
8. Locate “Network Security: LAN Manager Authentication Level” in the
list and double-click it.\
9. Change the setting from “Send NTMLv2 response only” to “Send LM &
NTLM - use NTLMv2 session if negotiated”\
10. Network Security: Minimum session security for NTLM SSP Based
(including secure RPC) Clients\
11. Change the setting from “require 128 bit” to unchecked (No Minimum)\
12. Click OK

the real difference between vista and windows 7 procedure is 10 and 11

It’s rather evil for Microsoft to choose defaults that make file sharing
not interoperable with latest Mac OS X. It’s not like they didn’t test
it so it must have been a conscious decision to screw with Mac users.
