Id: 944
Title: Short tutorial on svn propset for svn:externals property
Tags: svn
Date: 2006-06-06T19:53:38-07:00
Format: Markdown
--------------
Subversion has a nice way of including the content of one repository
into another repository. This is useful e.g. when you have a repository
with common routines that you want to use in multiple projects but you
don't want to duplicate the common code in multiple repositories (for
obvious reasons).

Using `svn:externals` property you can tell subversion to fetch content
of external repository in a given directory. However, the documentation
of `svn:externals` is weak so this post hopefully will save someone time
googling and figuring out how it works.

You set `svn:externals` property on an existing directory. The value of
this property is a list of directory/repository path values, separated
by spaces.

The best way to explain this is using an example. Usually you would use
just one external repository, but I'll use two in the example, just to
show how to create a list.

Let's assume that you have external subversion repositories
`http://svn.my.com/path/to/repo_one` and
`http://svn.my.com/path/to/repo_two`

You want `repo_one` in directory `dir_of_repo_one` and `repo_two` in
`dir_of_repo_two.` Create a text file with the value of `svn:externals`
property. Each list item is in it's own line, dir name is separated from
repo path by whitespace e.g.:\

    $ cat >svn_ext_val.txt
    dir_of_repo_one http://svn.my.com/path/to/repo_one
    dir_of_repo_two http://svn.my.com/path/to/repo_two

Now set the property on any directory already in subversion. In the
example it's the current directory:

    svn propset svn:externals . -F svn_ext_val.txt

Now when you do `svn update`, `dir_of_repo_one` will be created with the
content of `repo_one` and `dir_of_repo_two` with the content of
`repo_two`.

Voila! Pretty simple but you wouldn't know that just by reading the
docs.
