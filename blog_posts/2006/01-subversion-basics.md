Id: 304
Title: Subversion basics
Tags: svn
Date: 2006-01-02T16:00:00-08:00
Format: Markdown
--------------
**Installing subversion under Apache**:

* make apache a member of svn group: `usermod -G apache,svn apache`
* give apache permissions to `/home/svn`: `chmod g+rwx /home/svn`

**Backup and restore a repo:**

* `svnadmin dump repo >repo_svn_dump; bzip2 -9 repo_svn_dump`
* or: `svnadmin dump repo | bzip2 -c -9 >repo_svn_dump.bz2`

**Recreate a repo:**

* `svnadmin create repo`
* `sudo -u svn bzip2 -d -d repo.bz2 | svnadmin load repo`
* modify `/etc/httpd/conf.d/repo-svn-authzfile`
* you might, but don't need to, restart apache with `/etc/rc.d/init.d/httpd restart`

This assumes Apache2 and a particular setup found on RedHat Enterprise.

**Mark conflict as resolved**: `svn resolved $file`

**Update to a version in the past**: `svn update -r {2004-04-04}`

**Create a tag:**

`svn copy https://sumatrapdf.googlecode.com/svn/trunk/ https://sumatrapdf.googlecode.com/svn/tags/release-0.7 -m "Tagging release 0.7."`
