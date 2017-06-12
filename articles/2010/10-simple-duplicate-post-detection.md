Id: 336001
Title: Simple duplicate post detection for your blog, forum or commenting software
Tags: programming
Date: 2010-10-26T21:25:39-07:00
Format: Markdown
--------------
Mistakes happen
---------------

I was reading [scripting.com(scripting.com)](http://scripting.com) the
other day and some article had 3 identical comments from the same
person.

A mistake, obviously.

What is less obvious, at least to the authors of that commenting system,
is that such mistakes can be easily prevented with a bit of programming.

I know it’s easy because I’ve implemented it in a few lines of python in
software running this blog and in my [forum
software](http://blog.kowalczyk.info/software/fofou/index.html). They
run on App Engine but the technique applies to any web development
platform.

Duplication detection, the easy way
-----------------------------------

The idea is simple: before inserting an article/post/comment into a
database, check if an entry with exactly the same content already
exists. In most storage systems (like SQL database) doing this against
large text column is difficult and slow. To make it easy and fast we can
calculate SHA1 hash of the text, store it as part of the data describing
the post and check for duplicate hash.

SHA1 keys are short. They are a perfect match for key-value stores. With
a proper index they’re also very fast to check for in a SQL database.

For added robustness you can trim whitespace from the beginning and end
of text before calculating the hash.

This method doesn’t prevent malicious people (it only takes changing one
character to change the hash) but it does fix the common problem of
people submitting the same content twice because due to network
connection problems they were not properly notified that the post has
been successfully submitted. It happens more often that you might think.

Using other hash function, like MD5, will work too. SHA1 has better
properties but if your programming language provides MD5 but not SHA1,
use MD5.
