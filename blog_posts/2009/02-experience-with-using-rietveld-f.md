Id: 2075
Title: Experience with using Rietveld for code reviews
Tags: programming
Date: 2009-02-21T00:48:11-08:00
Format: Markdown
--------------
In my company we wanted to do code reviews. While code review process is
mostly about people, a good tool goes a long way.

My first thought was to use [ReviewBoard](http://www.review-board.org/),
but I was dreading installation process. Then I found out that
[Rietveld](http://code.google.com/p/rietveld/) can be easily installed
as an add-on for Google Apps domain.

The short version is: Rietveld works well. There were few gotchas:

-   Rietveld doesn’t support configuring with svn+ssh repositories which
    was important for us since our repository is not public. That means
    that during configuration I had to use a dummy repository link and
    that the only way to submit code for review is upload.py script
-   I have downloaded the public version of upload.py script. Big
    mistake: it doesn’t work with the version hosted with Google Apps.
    You have to download the script from your own instance
-   upload.py script often fails with 302 HTTP error. When that happens,
    delete cached authentication: `rm ~/.codereview_upload_cookies`

