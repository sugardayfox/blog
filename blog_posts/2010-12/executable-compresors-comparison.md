Id: 374001
Title: Executable compressors comparisons: upx 3.07w vs. mpress 2.17
Tags: note,sumatra,programming
Date: 2010-12-19T18:17:15-08:00
Format: Markdown
--------------
SumatraPDF executable is compressed with executable compressor. It makes
the file smaller and faster to download. On the downside, it increases
probability of false positive from virus scanners.

For a long time I’ve been using [upx](http://upx.sourceforge.net/) but
I’ll probably switch to [mpress](http://www.matcode.com/mpress.htm)
since it compresses a little bit better and it’s most expensive
compression `-s` is much faster than upx’s `—ultra-brute`.

Here’s a comparison using upx 3.07w and mpress 2.17 when compressing
release version of r2466 (pre-release for 1.3):

```
4598272 uncompressed
4527104 stripreloc /b /c

2205184 upx —best —compress-icons=0h
1756160 upx —ultra-brute —compress-icons=0

1725440 mpress -s (-s : find the best compression)
1735168 mpress -s -r (-r : dont compress resources)
1767424 mpress
```

The last result is for [StripReloc](http://www.jrsoftware.org/striprlc.php),
which is not a compressor. It only removes unneeded .reloc section in exes. Mpress must
already be doing that internally because compressing raw version and
stripreloc’ed resulted in the same size.

There is a serious downside to executable compressors: many anti-virus
program falsely report compressed executables as a virus or malware.
