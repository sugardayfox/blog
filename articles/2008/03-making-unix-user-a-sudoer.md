Id: 1010
Title: making unix user a sudoer
Tags: unix
Date: 2008-03-27T17:02:30-07:00
Format: Markdown
--------------
To make someone a sudoer:

Run `visudo` (must be root), edit `/etc/sudoer` and add:

```
%username ALL=(ALL) NOPASSWD: ALL
```

e.g.:
```
%kkowalczyk ALL=(ALL) NOPASSWD: ALL
```
