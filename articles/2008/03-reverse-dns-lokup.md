Id: 1231
Title: Reverse DNS lookup
Tags: unix,networking
Date: 2008-03-20T11:51:20-07:00
Format: Markdown
--------------
    From: http://www.cocoadev.com/index.pl?GettingHostnameFromIPAddress

    dig -x 17.149.160.49

    Listing all conections and who holds them:
    netstat -p -n | grep tcp | wc -l
    -p : show pid/name of the process
    -n : don't resolve name (much faster)

