Id: 1901
Title: Deep indentation vs. flat
Date: 2005-07-10T13:59:47-07:00
Format: Markdown
--------------
Having looked at lots of other people's code I've always felt alone in my
tendency to not use deeply nested indentation in my C code. Most code I've
seen looks like:


```java
if (foo) {
  if (bar) {
    if (gloo) {
    }
  } else {
  }
}
```

That's the "deep indentation" style. This is really reduced example and
doesn't show that in real code the else clause is usually empty or very short
and the code implements the pattern (if something succeeded and another thing
succeeded, and third thing succeeded, then do something useful), otherwise just
exit or return error code. I tend to write this type of code as:

```java
if (!foo)
  return;
if (!bar)
  return;
if (!gloo)
  return;
... and now do the thing
```

I find it easier to read and that's an important property of code. I always
thought that I'm the only one who does that, but now I know that [other people
do that too][1].

   [1]: http://wilshipley.com/blog/2005/07/code-insults-mark-i.html
