Id: 953
Title: Paradox of bad comments
Date: 2006-08-16T13:11:46-07:00
Format: Markdown
--------------
Bad comments are worse than no having comments at all. Given that
writing comments (good or bad) takes time, you would think that obviously bad
comments would be very rare.

In my experience, they’re not. Hence a paradox of bad comments.

I often see comments stating blindingly obvious things i.e. comments of
the kind:

```c++
class Foo {
 // constructor for Foo
 Foo();
};
```

or:

```c++
/* returns width */
int getWidth();
</code>

So why do people write such useless comments if they could get their job
done more quickly without spending time to write them?

My theory: guilt.

It’s not that those programmers don’t know that useless comments are,
well, useless, or that they couldn’t, given enough time to reflect, classify
such comments as useless.

Programmers know that writing good comments is important. However,
writing good comments is hard. By nature, good comments only explain tricky,
unexpected behavior of the code and those things are hard to explain
well.

On top of that, writing comments often has to be postponed until code
has been written and tested at which point there’s little incentive to add them.

Writing good comments is hard (which is why they’re rarely written) but\
programmers feel guilty when programs have no comments at all, so they
kill that guilty feeling by writing the easy, but useless, comments.
