Id: 2010
Title: Previewing changes before commiting on mac (svn or git)
Tags: git,svn,mac
Date: 2009-02-17T03:11:17-08:00
Format: Markdown
--------------
It’s a good habit to preview your changes before committing them. Long
time ago I wrote a program to help with that on Windows, called
[scdiff](/software/scdiff/). On mac, if you have TextMate installed, a
simple shell script does the job just as well:

```bash
#!/bin/sh

if [ -d “.svn” ]
then
 svn diff | mate
else
 git diff | mate
fi
```
