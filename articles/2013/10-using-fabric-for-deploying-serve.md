Id: 1416002
Title: Using Fabric for deploying server software
Tags: programming
Date: 2013-10-30T18:34:49-07:00
Format: Markdown
--------------
## Fabric - a tool for writing deployment scripts

Deploying a new version of server-side software is not trivial. Here are a few things that I need to do while deploying a new version of [App Translator](http://www.apptranslator.org) (written in Go) from my Mac dev computer to a production server:

 * compile the source .go files to an executable, on the server
 * which means I have to copy the source files from my dev machine to a server
 * ideally run tests to make sure I didn't introduce regressions
 * shut-down the currently running version
 * start the new version

Doing this by hand every time would be slow and tedious. Iteration speed [matters](http://www.azarask.in/blog/post/the-wrong-problem/).

The solution is automation, but that's tricky since it involves operations both on the local machine (like prepare the source for upload) and remote server (like compiling the sources on the server).

Luckily, I'm not the first person to encounter this problem which is why [Fabric](http://docs.fabfile.org/) exists.

Fabric is both a tool written in Python and a library of Python code for writing automatic deployment scripts.

## Installation

On my Mac I've installed Fabric with `sudo easy_install pip; sudo pip install fabric`

## Usage basics

The basics of Fabric are:

 * write the deploy script called `fabfile.py`
 * each function in `fabfile.py` can be thought of as an action
 * to execute a given action in `fabfile.py` do: `fab action`

As an example, I only have one action named `deploy` so to deploy a new version, I run `fab deploy`.

## Remote execution is via ssh

The code is deployed to Unix-based servers via ssh, so for convenience you can set up password-less (public key based) login for the account you use for deployment.

The code running locally runs under your account.

You need to specify which account on the server to use.

You can run an action against a single server or multiple servers at once.

## Possibilites

Fabric can run arbitrary code both locally and on the server, so the possibilities are, literally, limitless.

My [App Translator deploy script](https://github.com/kjk/apptranslator/blob/master/fabfile.py) is an example of a fairly complex script.

What you can learn from it:

 * how to specify server and user account for deployment (`env.user`, `env.hosts`
 * how to create a zip file locally, upload it to the server and unzip it there
 * how to compile sources both locally (no point uploading if the code doesn't compile) and on the server (to build the final executable)
 * how to run tests both locally and remotely
 * how to safely deploy a new version while preserving a way to easily revert to previous version in case of problems
 * how to do the the first-time setup (in case of deploying to a new server)

I've used Fabric in 3 of my projects and it does what it supposed to do well.

If I was writing new project from scratch, I would automate essentially everything related to setting up a new instance (it doesn't make sense to do it for existing projects since they're already set up).
