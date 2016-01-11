Id: 12012
Title: Compacting s3 aws logs
Tags: aws,python
Date: 2009-03-07T18:39:07-08:00
Format: Markdown
--------------
Amazon’s s3 can serve the files via HTTP protocol, so it can be used as
an HTTP server for static files. With a little effort you can configure
a given s3 bucket to save access log files to a different s3 bucket.
This is useful if you’re interested e.g. in figuring out how often a
given file is being downloaded. Logging doesn’t happen automatically. To
set it up, carefully follow [this official
document](http://docs.amazonwebservices.com/AmazonS3/latest/index.html?ServerLogs.html)

Unfortunately, standard s3 logs have a problem: they’re split among lots
of files. Around 200 log files are generated for one day which leads to
enormous number of files stored on s3 (e.g. merely 6 months of logs
generated \~35k of log files).

To fix that I wrote a python script that compacts the logs so that we
only get one file per day, which is a more reasonable number (35k files
shrank to only 171 for me). As a bonus, it also compresses the log.
Given very repetitive nature of log files, it’s a big win (the file
shrinks by 94% when compressed with bzip2 algorithm).

The logic of the script is simple:

-   download all the files for a given day locally
-   create one file that is compressed concatenation of those files
-   upload compressed file back to s3
-   delete the original files
-   repeat

It took me a few hours to get right, so in the spirit of sharing,
[here it is](http://code.google.com/p/kjk/source/browse/trunk/scripts/compact-s3-logs.py).
You’ll have to change a few things in the source:

-   s3 bucket name that stores the logs
-   where the logs are downloaded locally. I preserve them although they
    could be deleted if you’re short on space

Also, the script expects a `awscreds.py` file with access/secret key
i.e.:

```python
    access = “access key”
    secret = “secret key”
```

One unexpected thing I noticed is that apparently Amazon sometimes saves
the log file with screwy permissions so that it can’t be read or even
fixed by setting new acl (both those operations return 403 forbidden).
The only operation that works is deleting the file from s3, so that’s
what I do (it’s only a log file, after all).

Another thing I noticed is that boto sometimes hangs when talking to s3
(exception on timeout would be much better). If that happens, you can
kill the script and restart it - it’s smart enough to not re-download
files and only deletes the original log files on s3 after compressed log
has been saved in s3.

Expected usage of the script is as a daily script followed by a script
that analyzes the logs. Due to boto hanging it might not work well in
such scenario, though.
