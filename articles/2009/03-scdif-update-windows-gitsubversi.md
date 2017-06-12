Id: 15004
Title: scdiff update (Windows git/subversion/cvs gui diff previewer)
Tags: software
Date: 2009-03-11T21:40:44-07:00
Format: Markdown
Deleted: yes
--------------
I’ve updated my little [scdiff](/software/scdiff/index.html) program and
added git support.

If you work with git, Subversion or CVS from command line, scdiff is a
handy program to preview your changes before committing them.

It works by creating before/after files in a temporary directory and
launching a diff viewer (WinDiff by default, but can be configured) to
show the diffs.

It auto-detects which scm is being used in a given directory, so you
launch it by simply typing `scdiff`.

It requires git/svn/cvs executable be in the path. There still is no
decent windows build of git, so I use cygwin to get git.
