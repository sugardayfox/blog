---
Id: 20
Title: How to install latest clang (5.0) on Ubuntu 16.04 (xenial) / WSL
Date: 2017-11-20T20:37:02-08:00
Format: Markdown
---

This article describes installing latest clang (llvm) on Ubuntu 16.04 (Xenial), which is also the default distro for [Windows Subsystem for Linux](https://msdn.microsoft.com/en-us/commandline/wsl/about) (WSL).

Run:
```
wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | sudo apt-key add -
sudo apt-add-repository "deb http://apt.llvm.org/xenial/ llvm-toolchain-xenial-5.0 main"
sudo apt-get update
sudo apt-get -y clang-5.0
```

`clang-5.0` is the name of the executable (and so is `lldb-5.0` etc.).

## How it works

Let's deconstruct those commands so that you know what's happening.

How does apt know what packages are available?

Apt queries package servers to get a list of available deb packages. Default Ubuntu installation knows about official Ubuntu servers but you can run your own server to provide additional packages.

List of servers is in `/etc/apt/sources.list`. Here's how it looks by default on Ubuntu 16.04:

```
$ cat /etc/apt/sources.list
deb http://archive.ubuntu.com/ubuntu/ xenial main restricted universe multiverse
deb http://archive.ubuntu.com/ubuntu/ xenial-updates main restricted universe multiverse
deb http://security.ubuntu.com/ubuntu/ xenial-security main restricted universe multiverse
```

Ubuntu only provides a relatively old clang 3.8. Luckily, Apple creates deb packages and maintains a server for all llvm/clang releases and most Ubuntu distros.

`sudo apt-add-repository "deb http://apt.llvm.org/xenial/ llvm-toolchain-xenial-5.0 main"` adds llvm's server for Ubuntu 16.04 to `/etc/apt/sources.list`.

`sudo apt-get update` downloads the latest list of packages from all servers, including the one we just added

For security, packages are signed with private keys. You need public key to verify package signature.

`wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | sudo apt-key add -` downloads llvm's server public key.

`sudo apt-get -y clang-5.0` installs newly available package `clang-5.0`.

Flag `-y` disables confirmation prompt.

## What if there's a newer version of clang or a different version of Ubuntu?

There is a new llvm/clang release every 6 months. What to do for newer version?

Visit https://apt.llvm.org/ and locate the equivalent of `deb http://apt.llvm.org/xenial/ llvm-toolchain-xenial-5.0 main` for desired combo of clang/Ubuntu and correspondingly update `apt-add-repository ...` line in the above instructions.

