Id: 71001
Title: Thumbs up for markItUp
Tags: javascript,web devel
Date: 2010-04-18T15:23:25-07:00
Format: Markdown
Deleted: yes
--------------
[markItUp](http://markitup.jaysalvat.com/home/) is a simple enhancement
to textarea that adds icons that help writing markup for popular markup
languages (markdown, textile, bbcode).

It’s not wysiwyg but it does help in case you forget the markup syntax,
which happens to me a lot. The forgetting was bothering me enough to
improve the edit area on this blog, hence experimenting with markItUp.

markItUp turned out to be a good solution:

-   easy to integrate
-   supports multiple markup languages
-   uses jQuery (which is good for me, because I already use jQuery
    here, but might not be good if you’re using a different toolkit)
-   does have support for preview via server-side converter scripts. I
    haven’t used it - I’ve build preview system using the same ideas in
    just a few lines of code

I also looked at [wmd-new](http://github.com/derobins/wmd), which is
used on [StackOverflow](http://www.stackoverflow.com). I didn’t use it
because it only supports markdown (my markup of choice is textile,
mostly because markdown’s support for source code fragments is madness).

On the plus side, wmd can do previews without server support (it uses
markdown to html converter implemented in javascript).

I’ve also looked at rich-text editing like tinymce and ckeditor, but
they are pretty big, don’t support code snippets well and have way too
many buttons (my ideal rich text editor would support only small set a
features, somewhere between
[etherpad](http://code.google.com/p/etherpad/) and Google’s Weave
editor).
