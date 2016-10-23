Id: 1467
Title: Good programming practices
Date: 2002-11-17T05:31:05-08:00
Format: Markdown
--------------
[OSA Foundation](http://www.osafoundation.org) is a project to follow.
They have a high-quality people on the board, as evidenced by [this
post](http://lists.osafoundation.org/pipermail/dev/2002-November/000195.html)
on good programming practices:

> Off the top of my head, here's a quick partial list of things that I
> think are important:
>
> -   Continuous releases and daily builds
> -   Rather than put a bug in the bug database, first try to fix it.
>     Don't code first and ask QA to find bugs.
> -   The engineer who wrote the code is mostly responsible for QA.
> -   No code is finished until the unit test is finished.
> -   Good design results from getting something going quickly, followed
>     by lots of quick iteration.
> -   The best architectures evolve for a long time.
> -   Code reviews are a regular part of development.
> -   Difficult parts of code are best done by a team of two people
>     working closely together.
> -   If you see a new bug while doing something else, drop everything
>     and jump on the chance to fix it. It just make your life easier
>     down the road.
> -   It is possible to eliminate all bugs from many subsystems and you
>     should use techniques that make it possible, like monte carlo
>     simulations and exhaustive checking code whenever possible.
> -   Complexity is your enemy -- the sooner you learn the smell of too
>     much complexity, the better.
> -   Simplicity is your friend. Don't fall for a really cool fancy
>     algorithm unless you know you need it.
> -   Buzzword technologies often become religious. The only thing you
>     should be religious about is happy customers and they don't care
>     about buzzwords.
> -   Don't build elaborate systems unless you know you need them.
> -   Research existing open source project to make sure the code you
>     need isn't already available.
> -   Never write a slow program.
> -   By the time your project is finished your competitor's project
>     will be way better than you think.
> -   Price is only one feature of software. It's often not the most
>     important. From the user's perspective, simple things should be
>     simple and complex things should be possible.
> -   Automate releases.
> -   Use source level debuggers.
> -   Use IDEs
> -   Don't write makefiles unless you have too.
> -   Measure then optimize.
> -   Always spend at least 10% of your time learning new things.
> -   Use multiple languages, programming environments, operating
>     systems and programming methodologies until you are comfortable
>     with them and learn to appreciate the benefits of different
>     approaches.
> -   Use the best tools, even if they are expensive.
> -   Optimize variable names for maximum clarity, not minimum typing.
> -   It's OK to make mistakes. The only thing that isn't OK is to not
>     learn from them.
> -   Don't leave code commented out unless it's clearly explained.
> -   Approach each piece of code you write like it's final finished
>     code. Don't leave loose ends to clean up later because you'll
>     never have time to clean them up.
> -   Never code an algorithm that you know won't work in the final
>     product either because it won't scale, run fast enough, or will
>     have a rare failure case.
> -   All successful products have a life longer than you can conceive.
>     Files and processors will grow by factors of 10,000 -- make sure
>     your design can accommodate change.
> -   Don't code the first algorithm that comes to mind. Try to examine
>     all possible approaches and choose the best. Get someone to review
>     your approach before coding.
> -   Don't put band aids on bad code, rewrite it. At first is seems
>     hard, but after you've done it awhile you'll find successive
>     rewrites go much faster than you thought.
> -   When you work on someone else's code, don't leave it in worse
>     shape than it came to you. Ask the original author to review your
>     change.
> -   No bug is impossible to fix.
> -   Get to be good friends with a programmer, designer, or writer
>     who's better than you.

