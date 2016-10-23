Id: 179001
Title: Bash programming basics
Tags: reference,unix,note
Date: 2010-05-11T10:39:16-07:00
Format: Markdown
Deleted: yes
--------------
**Script arguments**

Script arguments are passed as `$0` (name of the script), `$1` (first
argument) etc.

`$*` - all arguments, as single word.\
`$`@ - all arguments, but each is seen as a separate word

**Check if a given argument was provided**

    if [ ! -n "$1" ]
    then
      echo "First argument not given"
    fi

**Handling errors from launched programs.**

Exit code from last executed command is in `$?`.

    if [ "$?" -ne "0" ]; then
      echo "last command failed"
      exit 1
    fi
