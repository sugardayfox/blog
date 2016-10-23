Id: 144001
Title: Idea for code review tool
Tags: idea,note
Date: 2010-05-04T11:06:24-07:00
Format: Markdown
--------------
Idea: code review tool for windows and mac.

Use cases:

-   preview your changes visually before submitting them to a
    repository. Enter commit message and commit from ui. Supports
    subversion and git.
-   ask for a code review by submitting files to website

Implementation plan:

-   Windows client
    -   get subversion files that need to be shown
    -   show list of files
    -   show a diff in ui
    -   allow entering commit message and do svn commit

Resources:

-   [diff library for c\# (Neil
    Frasner)](http://code.google.com/p/google-diff-match-patch/)
-   http://bazaar.launchpad.net/\~bzr-pqm/bzr/bzr.dev/annotate/head:/bzrlib/\_patiencediff\_c.c
-   http://bazaar.launchpad.net/\~bzr-pqm/bzr/bzr.dev/annotate/head:/bzrlib/\_patiencediff\_py.py
-   http://www.mathertel.de/Diff/default.aspx
-   http://www.menees.com/index.html
-   http://www.codeproject.com/KB/recipes/diffengine.aspx
-   http://www.eqqon.com/index.php/GitSharp
-   http://www.kaleidoscopeapp.com/
-   http://diffuse.sourceforge.net/ - in Python
-   http://winmerge.org
-   [Review Board](http://www.reviewboard.org/)
-   [Rietveld](http://code.google.com/p/rietveld/)
-   http://stackoverflow.com/questions/145607/text-difference-algorithm

Things todo for korev:

-   parse output of git status, get added/changed files
-   show list of added/deleted/changed files (if more than one)
-   use git to get before/after version of changed files
-   show added files as-is in a text viewer
-   show changed files as a diff (using DiffMatchPatch.cs to get diff
    info)
-   detect that git is not installed (try to execute git —version) and
    show error message

Ways to tackle this:

-   write a simple text viewer with the simplest way of spliting files\
     into lines and without a way to wrap)
-   simple way as above but with support for wrapping
-   write a simple text viewer by porting qemacs way of storing data
-   re-use editor code from SharpDevelop

`git show master:win/korev/App.xaml`

`git show master~1:win/korev/App.xaml`

git:

-   working tree - what’s on disk
-   index - what has been added with git add, the content at the time of
    adding
-   repository - sutff that has been git commit’ed

Git commands:

-   `git diff` - changes in working tree that are not yet staged for
    next commit
-   `git diff --cached` - changes between index and last commit, what
    would be commited with git commit
-   `git diff HEAD` - changes in working tree since last commit, what
    would be commited with git commit -a
-   `git diff HEAD^ HEAD` - compare version before last commit with last
    commit

`git log --raw` : get parseable history of changes\
`git status --porcelain` : get parseable status

Related idea:

-   a website for browsing the source code of a project. Focused on
    studying one version of the code
-   a website for browsing history of changes in the source code.
    Focused on figuring out what changed when
-   a website for targeted searched of the code. Like google code search
    but more targeted i.e. where you first select a language (C\#, Java,
    Python etc.) and has a curated list of source code repositories to
    search, not “everything goes” like google code search
-   port Review Board from sql to redis, to see how faster it would be

