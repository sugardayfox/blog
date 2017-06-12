Id: 1853
Title: Don't use 0 instead of NULL
Date: 2004-07-22T10:20:45-07:00
Format: Markdown
--------------
Linus [has spoken](http://lwn.net/Articles/93577/): using 0 (instead of
NULL) to denote a null pointer is wrong:\

    char * p = 0;    /* IS WRONG! DAMMIT! */
    int i = NULL;  /* THIS IS WRONG TOO! */

\
I concur.
