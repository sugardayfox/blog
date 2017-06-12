Id: 8001
Title: Where do bugs come from?
Tags: programming
Date: 2009-02-24T09:59:29-08:00
Format: Markdown
--------------
Where do bugs come from?
------------------------

Let’s assume you know, more or less, what you want to write and are
ready to start cranking the code. You’ll write bugs. What’s the source
of those bugs and how can you avoid them?

Ignorance and carelessness
--------------------------

There are only 2 sources of bugs in software:

-   ignorance
-   carelessness

Many faces of ignorance
-----------------------

Ignorance is not knowing what you’re doing. It’s a broad category, so
let’s list some of the many ways to be ignorant.

### Not knowing the language

For example, you don’t know that in C you need to free() each malloc().
Or that you can’t return a pointer to a local variable. Or you don’t
understand the difference between sizeof(TCHAR) and sizeof(TCHAR\*). Or
any other of the infinite number of things to know. It takes years to
master a programming language.

### Not knowing the language very well

Even if you’re not writing outright bugs, there are still millions of
ways to misuse a language. A trivial example in C might be malloc()ing a
10-byte temporary buffer in a function instead of using stack (which is
faster and generates smaller code). Again, it takes years to really
master a language.

### Not knowing operating system and its APIs.

Windows has thousands of APIs. Mac OS has thousands of APIs. Way too
much to keep details of every API call in your head. It’s easy to misuse
the APIs and misuse leads to bugs.

### Not knowing file format/protocol/specification.

Even if you know your language and APIs, there’s still plenty of sources
of potential ignorance. Software doesn’t live in a vacuum but in a very
rich and complicated ecosystem of existing solutions that are usually
poorly designed or over-complicated due to historical reasons that no
longer apply. Unfortunately our software has no choice but to talk to
those systems.

If your software talks HTTP protocol, did you really carefully read the
full HTTP 1.1 spec and made sure to handle all possible cases?

If you use XML, did you take time to read and fully comprehend XML spec?

If you’re writing a text editor, did you research character encodings,
Unicode, different line endings on various platforms?

It’s easy to gain a superficial understanding of a complex system. Isn’t
XML just stuff inside brackets? Bugs come from missed details. In XML’s
case it’s handling BOM, all possible encodings of XML files, entities,
CDATA sections, comments, knowing which bytes cannot be put inside XML
file (and more).

Inevitable mistakes
-------------------

Even if we assume that we know everything, there’s the issue of
mistakes. Even an experienced writer cannot write few pages of text
without making a grammatical mistake. Programming is no different - no
matter how good you are, you’ll make a number of mistakes per N lines of
code.

The more careful you are, the less mistakes you will make, but no matter
what the effort, the mistakes will be there

Avoiding bugs - the art of education and compensation
-----------------------------------------------------

Now that we know where bugs come from, what can we do to decrease the
number of bugs?

The methods (and apparatuses) we can use to decrease number of bugs in
software fall into two categories:

-   education
-   compensation

We educate ourselves to minimize ignorance. We compensate for inherent
fallibility of human mind to reduce number of mistakes we make or find
them faster. Here’s a list of some specific things you can do.

### Read

Read books, papers and web articles about your language, operating
system, tools etc. That’s pretty obvious.

Study the APIs you use. Way too often we use APIs having only partial
understanding of what they do, how can they fail, what consequences do
they have etc. For Windows, MSDN documents API calls in great detail.
Read the whole description, including remarks, not just the API
signature. A similar resource to MSDN exists for every system out there.

Study the source code of other software. This is an open-source era. We
have source code for thousands of systems, small and large, easily
available for us to learn from. There’s no excuse not to do it.

You can read the source code in an exploratory way i.e. randomly browse
the code fishing for nice tricks you can adopt in your own software.

You can do it in a focused way e.g. if you know you need to implement
hash table, you can look at source of different implementations of hash
tables to learn the techniques used.

Look for low-level implementation tricks and techniques as well as big
picture (e.g. how the code is organized into modules).

Reading code is hard and boring. Even high-level languages operate on a
very low level and it’s hard to see the forest because all those trees
obstruct the view. It is, however, a skill that can be improved with
practice.

Be reflective. When you find it hard to understand the code, try to
figure out why is it hard (lack of comments? confusing comments? poor
structure? poor naming of variables and functions?) and avoid doing
things that make the code hard to read in your own code.

Conversely, when you see exceptionally well written code, try to figure
out what makes it good and imitate.

### Ask for code reviews

The amount of knowledge about programming is enormous. It takes years to
master just one programming language or API set. Ignorance is a
permanent state of being for a programmer and the worst part is that we
don’t know what we don’t know.

Wouldn’t it be great if there was a magic fairy that would just tell us
about the mistakes we made? No such fairy exists but the next best thing
is your coworker. He might know that you should free() your malloc()s.
He might know a better (simpler, faster) way to do something. He might
spot a mistake. All you need to do is ask for a code review.

### Talk about programming

Discuss programming topics with your peers. Ask them to teach you tricks
they’ve learned. Tell them about the tricks you’ve learned. No matter
how much you know, there’s always more to learn. For example, after
umpteen years of writing code I’ve learned a cool linked list trick from
interview with Anders Hejlsberg (it’s at the very end of the interview).

### Use automatic tools for finding bugs.

Static checker analyzes the source code looking for mistakes. There are
static code checkers for Java (FindBugs), C\# (FxCop), Python (PyLint),
C/C++ (Coverity, Visual Studio 2005 Team (aka. “expensive”) Edition has
integrated static checker) and probably other languages.

There are run time tools like BoundsChecker, AppVerifier, Valgrind,
AQTime and a countless number of memory debugging tools.

Learn about those tools and start using them regularly.

### Use assert()

In C assert() macro checks if an expression is true. If it’s not, it’ll
do something to alert the programmer (details depend on implementation,
in gcc it calls abort(), in Visual Studio it shows a dialog box that
gives the programmer to break into an app and inspect the state of the
program).

Assert-like functionality exists in other languages as well.

The idea is to add checks for situations that should never, ever happen.

There’s a subtle distinction between something that almost never happens
but is possible and something that should never, ever happen and I’ve
seen asserts misused because people fail to understand that distinction
(again, ignorance!).

The most common example of misusing asserts is validating that malloc()
didn’t return NULL. While it’s unlikely to happen on a desktop machine,
it’s a perfectly valid for malloc() to fail and return NULL.

An example of a valid assert is checking that after you insert a node in
a list, the length of the list will increase by one if the insert was
successful.

Assert impossible, not unlikely.

I’ve learned to assert() a lot while working at Microsoft since the code
I worked on was littered with asserts (an example of learning good
habits by reading other people’s code). Many times asserts noticed me
about bugs that I would have missed and that’s why I think it’s
worthwhile to use them.

### Write unit tests.

Lots have been said on this topic since eXtreme programming became
popular. I’m not going to explain why it’s a good thing - others have
done it better that I could.

The reality of unit-testing is not as good as XP hype would have you
believed. My random sampling of open-source projects and projects I’ve
seen inside companies shows that almost no-one writes unit tests

Writing unit tests is difficult because it’s boring. When we understand
the problem, our first instinct is to dive in and code the damn thing
because it’s the shortest path to having things work. At the time of
writing the code, it has obvious cost (time spent writing additional
code for unit tests) but no obvious benefit.

Unit tests do payoff in longer term when they save you from long,
hair-pulling debugging sessions. Tthe longer you expect the code to
live, the bigger payoff from writing unit tests.

Not that everyone agrees with that

### Write automated tests

Manual testing i.e. running your software and using it to see if it
works as expected, is expensive. Automated testing is also expensive to
setup and maintain, but it pays off in the long run.

Unit tests is only one example of possible automated tests. Unit tests
are usually confined to thoroughly testing a small thing: a function or
a class. There are other ways to automatically test the code that are
beneficial. The golden rule of testing is that the more of it you have,
the more bugs it’ll catch, the better the quality of your software.

For example, in Mozilla project, Java Script engine has a large number
of relatively small Java Script code snippets that ensure that the
implementation works the way it was designed. Also in Mozilla, Gecko
rendering engine has has a number of layout test cases ensuring that
changes to layout engine code do not change the way HTML pages are
rendered on the screen. Sqlite database engine has a battery of tests to
exercise the engine.

Automated tests are a rich topic. Think about ways you can automatically
test your code (different strategies work for different kinds of code).

### Do manual testing.

Oh yeah, just because you have automated tests doesn’t mean you can skip
manual testing. Let me remind you a golden rule of testing: the more of
it, the better.
