Id: 14
Title: Go package for better guid generation
Date: 2015-02-12T17:17:01-08:00
Tags: go, programming
Format: Markdown
--------------
The need to generate a globally unique identifier comes up often. The way described in [RFC 4122](https://tools.ietf.org/html/rfc4122) is popular but it can be done better.
I wrote [betterguid](https://github.com/kjk/betterguid) Go package that does it better. The properties of this method are:

* generated id is 20 characters, safe to include in urls (no need to escape)
* they are based on timestamp so that they sort **after** any existing ids
* they contain 72-bits of random data after the timestamp so that IDs won't collide with other clients' IDs
* they sort **lexicographically** (so the timestamp is converted to characters that will sort properly)
* they're monotonically increasing. Even if you generate more than one in the same timestamp, thelatter ones will sort after the former ones. We do this by using the previous random bits but "incrementing" them by 1 (only in the case of a timestamp collision).

You can read a [longer description](https://www.firebase.com/blog/2015-02-11-firebase-unique-identifiers.html) of the algorithm.
My implementation is based on this [JavaScript code](https://gist.github.com/mikelehen/3596a30bd69384624c11).
