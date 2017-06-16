Id: 238002
Title: Hiding duplicate content from your site via robots.txt
Date: 2012-10-22T01:43:33-07:00
Format: Markdown
--------------
Many blogs, including this one, generate duplicate content. For example,
[the archive pages](/archives.html) duplicate the
content of individual posts, they just show them in a different way (a
couple of posts per page, as opposed to a single post per page).

That unfortunately clogs search engines. Being a perfectionist that I
am, I want that a search for e.g. “15minutes” (my [simple timer
application](/software/15minutes/)) leads people to individual blog
posts about it and not to aggregate pages with random other content.

Thankfully there’s a way to tell search engines to not index parts of
your site. It’s [quite
simple](http://www.javascriptkit.com/howto/robots.shtml) and in five
minutes I cooked up the following [robots.txt](/robots.txt) for my
site:

```
User-agent: *
Disallow: /page/
Disallow: /tag/
Disallow: /notes/
```

In my particular case, archive pages all start with `/page/` or
`/notes/` and `/tag/` is another namespace with duplicate content (shows
a list of articles with a given tag).

For this technique to work the names duplicate pages have to follow a
pattern, but that’s easy enough to ensure, especially if you write your
own blog software, like [I do](http://github.com/kjk/web-blog).
