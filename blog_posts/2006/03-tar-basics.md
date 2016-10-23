Id: 1291
Title: tar basics
Tags: unix
Date: 2006-03-30T16:00:00-08:00
Format: Markdown
--------------
**Unpack an archive**:\
`tar xvf archive.tgz`

**Create an archive**:\
`tar zcvf archive.tgz archive-dir`

**Unpack in a given directory**:\
`tar -C dir -zxvf name.tgz`

Note: dir must already exist

**Options:**

    -v     : verbose
    -x     : extract
    -c     : create
    -f     : file name
    -C dir : where to uncompress

Options can be mushed (i.e. `-xvf` is same as `-x -v -f`)
