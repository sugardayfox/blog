Id: 919
Title: Unsolved source control problems
Date: 2005-10-24T21:06:19-07:00
Format: Markdown
--------------
Recently there was an explosion of new source control programs:
subversion, arch (and its many children), darcs, monotone, git,
codeville, mercurial (and more).

One of the reasons for creating them was the fact that CVS, the king of
the hill so far, lacks some basic features. Or, as a less cultured
person might say, sucks.

Almost all those new programs also try to make open-source, distributed
collaboration easier. Ironically, I don’t see any of them solving the
problems that I’ve come across trying to make small contribution to
various open-source projects so I don’t think that they are solving the
right problems.

In open-source projects there are 2 sides:

-   project leader (or leaders) - those are the people who can commit
    changes\
    to the repository
-   contributors - a casual contributor didn’t earn enough street cred
    to get\
    full write access

First problem I saw is that submitting patches, reviewing and accepting
them is time consuming for both leaders and contributors which causes
projects to suffer because contributions are ignored which discourages
contributors from contributing more.

Ideally, the flow would be:

-   contributor does some code modifications
-   contributor notifies the leader
-   leader looks at modifications
-   leader applies modification to the main tree or marks them as
    rejected,\
    providing explanation why they were rejected

With CVS it’s all very painful. Contributor has to generate patches. He
has to send them via e-mail or attach to a bug tracking system. Leader
has to apply the patch to his tree (which might be out-of-sync wrt. to
the tree against the patch was generated), preview the patch, respond
via e-mail or by making a note in the bug tracker.

I went through this a couple of times and for small changes it’s just
not worth the effort, especially if patches rot in the bug tracker for
weeks. But if a small change is not accepted quickly, I’m not going to
bother spending more time developing more extensive enhancements,
therefore the project looses me as a potential contributor.

Now, it might be just a social problem (it’s hard to imagine a tool that
would force project leaders to respond to submitted patches faster) but
I believe that a better tool would help a lot.

One option would be to give repository write access to anyone who asks,
but the cure could be worse than a disease.

Some systems (e.g. darcs, arch) say that they’ll solve this problem by
making things more decentralized i.e. by making it easier to create your
own repositories based on the main tree and making it easier to merge
changes between different repositories. But this does have a few
problems of its own:

-   not everyone can make their repositories available to other people
    (you need at least a web server)
-   I think that having one repository is better; another repository is
    a mini-fork and there’s no easy way to tell other people about it

I think that all a source control system needs to make things better is:

-   anonymous users should be able to create a branch
-   branch management should be easy (list all branches, preview changes
    in a given branch, merge a change from one branch to another)
-   ability to comment on changes

That way the work flow could be:

-   contributor creates a branch for his changes
-   contributor notifies the leader e.g. by e-mail “hey, I’ve made some
    changes in a branch foo”
-   leader can issue one command that would show him the changes made on
    branch foo, as a diff or in some gui
-   if leader likes the changes, he can merge them back to main tree
    with one command; if he doesn’t he can add a comment explaining why
    changes are not good
-   if contributor sees the comment, he can modify the code according to
    comments and redo the process

It’s almost like the full decentralized model but separate repositories
would be replaced by branches in the main repository. Anyone could
create a branch and make his changes on his own branch. Having
everything in one place would make it easier for other contributors to
discover interesting changes.

Currently the biggest problem with open-source development I see is that
whoever leads the project has a strong-hold on what goes in and what’s
not. Having a system that supports branches for everyone but maintains a
restricted access to the mainline would hopefully encourage more
experimentation, more people contributing to projects.

Which reminds me that I’m very disappointed by SourceForge. Did you
notice that it hasn’t changed in ages? No new cool Web 2.0 AJAX-y
interfaces.

It still only offers CVS. They don’t offer built-in wiki for projects.
Their implementation of a mailing list and a bug tracker is terrible.
And, ironically, the source code to the biggest open-source code
repository in the world is closed.

SourceForge could be so much more…
