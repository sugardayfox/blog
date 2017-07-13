---
Id: wpig-02
Title: Installing necessary tools
Format: Markdown
Tags: go
CreatedAt: 2017-07-12
PublishedOn: 2017-07-12
Collection: go-windows
Description: Installing Go toolchain and necessary libraries.
---

You probably have the tools already installed, but I have to cover it anyway. Such is life of a tutorial writer.

First, you need to install Go compiler. Go to [official download page](https://golang.org/dl/). When I write this the latest version of Go compiler is 1.8.3 so that's what I'll use as example. You should adapt the instructions to the latest version.

There are 2 versions of the compiler:
* 32-bit version (go1.8.3.windows-386.msi)
* 64-bit version (go1.8.3.windows-amd64.msi)

32-bit version of the toolchain can run on both 32-bit and 64-bit version of Windows. By default it produces 32-bit executables but can also generate 64-bit executables.

64-bit version of the toolchan can only run on 64-bit version of Windows. By default it produces 64-bit executables but can also produce 32-bit.

After installing Go toolchain, you need to setup `$GOHOME` environment variable. I set it to `$HOME\src\go`.

**TODO:** describe how to do it. Maybe screenshot of env variable editor or describe how to setup powershell using `$profile`.

This book heavily uses 2 libraries:
* https://github.com/lxn/win provides Go bindings to win32 API calls
* https://github.com/lxn/walk builds on top of `lxn/win` and provides higher-level wrappers and declarative UI and layout framework. Think of it as equivalent of winforms in C#

You might just as well install them right away by doing `go get github.com/lxn/win` and `go get github.com/lxn/walk`.

I'm not one to tell you which text editor to use. I use [Visual Studio Code](https://code.visualstudio.com/) with [Go extension](https://marketplace.visualstudio.com/items?itemName=lukehoban.Go) and it works very well.
