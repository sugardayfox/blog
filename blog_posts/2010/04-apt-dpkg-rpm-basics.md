Id: 115001
Title: apt, dpkg, rpm basics
Tags: unix,note
Date: 2010-04-29T16:02:20-07:00
Format: Markdown
Deleted: yes
--------------
**apt basics:**

    apt-get install ${package}      : install a package
    dpkg --list                     : list installed packages
    apt-get remove ${package}       : remove a package
    apt-cache search ${pattern}     : find packages matching pattern
    apt-cache show ${package}       : show info about a package
    sudo apt-get install apt-file   : apt-file needs to be installed first
    apt-file show ${package}        : shows files in a package
    apt-get upgrade                 : upgrade all installed packages

Note: `apt-file` might need to be installed first
(`sudo apt-get install apt-file`)

**dpkg basics:**

    dpkg -i foo.deb      : install a package
    dpkg -l | grep foo   : list installed packages, show those matching foo
    dpkg -r foo-ver      : remove pacakge

**rpm basics:**

    rpm -ihv *.rpm         : install a package
    rpm -Uhv *.rpm         : upgrade a package
    rpm -qa | grep -i name : list all packages matching a name
    rpm -ql $package-name  : list files in a given pacakge
    rpm -q --whatrequires mysql-server : list packages dependent on a given package
    rpm --recompile $package-name      : build and install a source package
