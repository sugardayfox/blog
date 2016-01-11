Id: 25011
Title: Network drives, .net, security and virtualbox
Tags: .net,c#
Date: 2009-06-04T23:28:28-07:00
Format: Markdown
--------------
I have Windows 7 installed on my MacBook under VirtualBox. To keep
things simple, I keep my code on mac side, sharing the directory with
Windows side.

Tip #1: VirtualBox 2.2.4 has weak implementation of shared folders. For
example, Visual Studio 2008 sees the files on shared drive as
write-protected (even though other applications don’t). Visual Studio
2010 doesn’t have this problem, but it can’t finish compilation,
complaining it cannot write some intermediary build files to a shared
folder. Also, trying to mark shared folder as trustworthy for .NET
doesn’t seem to work (see below).

Solution: use Mac’s built-in file sharing instead and shared things out
using built-in samba (smb) sharing in System Preferences/Sharing.

Another problem is that .NET by default doesn’t trust code on network
drives and Visual Studio will tell you that “The Project Location is Not
Trusted” if you open a project from a network drive. The magic
incantation to make a given shared folder (in my case it’s drive z:)
trusted is:

```
caspol -q -machine -addgroup 1.2 -url file://z:/* FullTrust
```

This command has to be executed as Administrator.

One way to have `caspol.exe` available in the PATH is to start Visual
Studio cmd-line (shortcut is installed in Start menu).
