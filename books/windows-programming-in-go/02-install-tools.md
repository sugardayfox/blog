---
Id: wpig-02
Title: Installing necessary tools
Format: Markdown
Tags: go
CreatedAt: 2017-07-12
PublishedOn: 2017-09-12
Collection: go-windows
Description: Installing Go toolchain and necessary libraries.
Status: draft
---

You probably have the tools already installed, but I have to cover it anyway. Such is life of a tutorial writer.

First, you need to install Go compiler. Go to [official download page](https://golang.org/dl/). When I write this the latest version of Go compiler is 1.8.3 so that's what I'll use as example. You should adapt the instructions to the latest version.

There are 2 versions of the compiler:
* 32-bit version (go1.8.3.windows-386.msi)
* 64-bit version (go1.8.3.windows-amd64.msi)

32-bit version of the toolchain can run on both 32-bit and 64-bit version of Windows. By default it produces 32-bit executables but can also generate 64-bit executables.

64-bit version of the toolchan can only run on 64-bit version of Windows. By default it produces 64-bit executables but can also produce 32-bit.

You can follow [official installation instructions](https://golang.org/doc/install) or the steps below.

After installing Go toolchain, you need to setup `$GOPATH` environment variable. I set it to `$USERPROFILE\src\go`.

You also need to add `$GOPATH\bin` to `$PATH`. When you `go get` a program (as opposed to just a library), go toolchain installs it in `$GOPATH\bin`. We want them to be avilable in `$PATH` for convenience. We'll need some tools installed via `go get` on this journey.

If you always use powershell, you can add the following to `$profile` file (e.g. via `notepad $profile`):
```
$env:GOPATH = $env:USERPROFILE + "\src\go"
$env:Path = $env:Path + ";" + $env::GOPATH + "\bin"
```

`$profile` is executed when powershell starts.

Don't forget to create `$env:GOPATH` directory.

**TODO:** describe how to do it. Maybe screenshot of env variable editor or describe how to setup powershell using `$profile`.

This book heavily uses 2 libraries:
* https://github.com/lxn/win provides Go bindings to win32 API calls
* https://github.com/lxn/walk builds on top of `lxn/win` and provides higher-level wrappers and declarative UI and layout framework. Think of it as equivalent of winforms in C#

Install
* `go get github.com/lxn/walk` : this is the library we'll be using
* `go get github.com/akavel/rsrc` : we'll need `rsrc` tool

I'm not one to tell you which text editor to use.

I use [Visual Studio Code](https://code.visualstudio.com/) with [Go extension](https://marketplace.visualstudio.com/items?itemName=lukehoban.Go) and it works very well.
