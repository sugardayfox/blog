Id: 1964
Title: Deeply nested if statements
Date: 2006-08-22T13:56:32-07:00
Format: Markdown
--------------
I dislike deeply nested C\C++ code. I.e. the code of the form:

```java
if (foo) {
  if (bar) {
   if (anotherVariable) {
   }
  }
}
```

It's a simplified example that doesn't really show the problem with such code
- when the logic within an if is long, such code becomes hard to read because
it becomes difficult to figure out what logic condition is handled by a given
part.


To my surprise I seem to be in minority. Most people, especially the Windows
GUI kind, nest logic conditions like there's no tomorrow. I find it hard to
argue about the superiority of avoiding nesting because saying "this is hard
to read" is too vague. But I'm relieved to find out that I'm not the only one
that has a dislike for unnecessary nesting. Stepanov, inventor and implementor
of STL, has a great example of simplifying nested logic in his
"Professionalism in Programming" papler ([PowerPoint][1], [PDF][2], from [his
website][3]).


The way to avoid deeply nesting is to do early exit as soon as possible. The
trivial example could be rewritten as:

```java
  if (!foo) return;
  if (!bar) return;
  if (!anotherVariable) return;

  ... and this is the logic
```

Of course trivial examples teach us nothing - this only shows the forest so
that you don't get lost in details. For a real-life example see how Stepanov
refactors deeply nested logic expression (pages 4-13 in the paper above).


   [1]: http://www.stepanovpapers.com/Professionalism%20in%20Programming.ppt

   [2]: http://www.stepanovpapers.com/Professionalism%20in%20Programming.pdf

   [3]: http://www.stepanovpapers.com/
