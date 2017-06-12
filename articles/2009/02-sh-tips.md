Id: 7051
Title: ssh tips
Tags: ssh,unix,svn,reference
Date: 2009-02-24T11:05:39-08:00
Format: Markdown
--------------
**Password-less authentication:**

-   generate public and private key (e.g. with
    `ssh-keygen -q -f ~/.ssh/id_rsa -t rsa`) (or `-t dsa`)
-   on the machine from which you connect: put private key in
    `~/.ssh/id_rsa` (for rss key) or `~/.ssh/id_dsa` (for dsa key), with
    0600 permissions
-   on the machine you connect to: append public key to
    `~/.ssh/authorized_keys` (also with 0600 permission)

You can read [long version](http://sial.org/howto/openssh/publickey-auth/)

**Using a different username with svn+ssh:**

By default ssh uses \$USER when connecting to a server. This makes
things difficult if you use svn with ssh authentication and your
username for svn server is different than \$USER. To get around this
define per host SSH settings. To do this, edit `~/.ssh/config` and add
the lines:

    Host svn.yourserver.com
    User jdoe

That instructs ssh to use the username jdoe instead of your unix
username when connecting to svn.youserver.com. If you use
password-protected ssh keys, you will be prompted for a private key
password, like this:

`Enter passphrase for key '/home/jdoe/.ssh/id_dsa':`

To get around this you can either setup ssh-agent to cache this password
or you can turn off the password on your private key. To do the latter:
`ssh-keygen -t dsa -p` and when it prompts for a new passphrase, just
hit enter.

**Using ssh agent:**

Ssh agent remembers password for your private key, so you only have to
type it once in a session, not every time it’s used.

Start ssh agent: **ssh-agent**\
Add keys: **ssh-add**\
List currently loaded keys: **ssh-add -l**

**Cygwin and ssh agent:**

Links:

-   <http://holdenweb.blogspot.com/2007/12/cygwin-ssh-agent-control.html>
-   <http://www.newmedialogic.com/node/55>

Add this to `~/.bashrc`:

```

## Enable ssh agent
export SSH\_AUTH\_SOCK=/tmp/.ssh-socket

ssh-add -l \>/dev/null 2\>&1

if [ $? = 2 ]; then
 # Exit status 2 means couldn’t connect to ssh-agent; start one now
 rm -rf /tmp/.ssh*
 ssh-agent -a $SSH_AUTH_SOCK >/tmp/.ssh-script
 . /tmp/.ssh-script
 echo $SSH_AGENT_PID >/tmp/.ssh-agent-pid
fi

function kill-agent {
 pid=`cat /tmp/.ssh-agent-pid`
 kill $pid
}

function addkeys {
 ssh-add ~/.ssh/id_dsa*
}
```

**Set an alias for a hostname**

In `~/.ssh/config`:

```
    Host dev30
    HostName dev30.kowalczyk.info
    User root
    IdentityFile ~/.ssh/id_dsa_dev30
```
