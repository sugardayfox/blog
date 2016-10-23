Id: 1308
Title: Local DNS modifications on Windows (/etc/hosts equivalent)
Tags: windows
Date: 2005-12-30T16:00:00-08:00
Format: Markdown
--------------
On Unix, `/etc/hosts` file contains mappings between an IP address and a name of the host. It overrides mappings from DNS. Windows has an equivalent of this file: `c:\WINDOWS\system32\drivers\etc\hosts` (at least that's the name on Windows XP).

What are possible uses for this file:

* when I make changes to my blog, I develop and test the changes on my local Apache setup on windows. My config files are setup up so that they work correctly both on my local test server and deployed server. To make things easier, I use hosts file to map local loopback address 127.0.0.1 to blog.local.org
* changing hosts file can be used to fool local computer and simulate DNS changes (e.g. for testing before making real DNS changes)

Changes to hosts file take place immediately after saving the file.

Here's the simplest mapping from 127.0.0.1 to localhost and blog.local.org names:
```
127.0.0.1       localhost blog.local.org
```
