Id: 1863
Title: scdiff 0.3 released
Date: 2004-10-02T21:43:39-07:00
Format: Markdown
--------------
I've released version 0.3 of [scdiff](/software/scdiff). The fixes in
this version:\

-   [Eli Tucker](http://nerdmonkey.com) fixed a bug handling directories
    in cvs\
-   no longer crash when handling `svn remove`. Still, handling of
    deleting files\
     is far from perfect (currently they are silently ignored while we
    should show apropriate diff)
-   properly handle file names with spaces. Fix suggested by Raman Gupta

