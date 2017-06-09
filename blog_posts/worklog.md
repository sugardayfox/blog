2017-06-09

#read http://floooh.github.io/2017/06/09/webassembly-demystified.html #webassembly

2017-06-08

#read http://www.kalzumeus.com/2008/01/28/why-you-shouldnt-pay-any-seo-you-can-afford/ #seo

---

#flexbox #css

To implement layout: [[left]     [right]]:

div style=display:flex;flex-direction: row;justify-content: space-between;
  div
    left
  div
    right

---

Sciter (https://sciter.com/) is a tool for writing desktop UI quickly. It's like Electron but much smaller (5-8 MB vs. hundrends of megabytes). It's almost HTML and almost JavaScript.

It's also under-documented. Here's how I started playing with it:
* downloaded SDK (https://sciter.com/download/) into ~/Downloads folder and unzip
* create alias for the test app: alias sciter='/Users/kjk/Downloads/sciter-sdk/bin.osx/scapp'
* added the above alias to ~/.bash_profile

When I'm prototyping:
* create app.html and write the code
* to run: sciter app.html

Given how poorly documented Sciter is, I mostly use samples/ from the SDK to learn how do interesting things.

Ultimately the idea is that you use sciter .dll as part of your native Windows (C++ or .NET) or Mac (Objective-C or Swift), with majority of the interface implemented in Sciter's dialect of HTML/JavaScript, exposing functionality from native code as tiscript code.

Using scapp is a way to quickly prototype UI.

#sciter

---

#go Implemented worklog for blog.

---

Reading more blog posts from http://thestartuptoolkit.com/.

Insight: if your business is "better X", you should use X to make sure you're really improving on it. http://thestartuptoolkit.com/blog/2011/09/use-the-tools-youre-displacing/

#business

---

https://blog.figma.com/webassembly-cut-figmas-load-time-by-3x-76f3f2395164. Summary: webassembly is much faster than asm.js. #webassembly

---

#postgres http://www.craigkerstiens.com/2017/06/08/working-with-time-in-postgres/

---
#idea

For faster PDF rendering, re-compile it. PDF format is rather complicated, mostly text format. A PDF renderer first has to parse it into in-memory representation and then renders it.

Imagine a binary format that is designed for fast loading and one-time PDF => binary format compilation step. With the right format, parsing step could essentially be eliminated.

We could even improve rendering time with some additional optimizations.

For example, there are PDF documents with gigantic images (think 4000 x 4000 pixels) that could be resized to some reasonable value and saved in a format that is fastest to decode.

I'm sure there are optimizations possible for complex vector graphics (e.g. removing stuff that we know will be invisible).
---

#watched https://vimeo.com/77265280, trying to understand capabilities of WebRTC. I want to build VNC client in the browser but without ability to open tcp/udp connections it's not possible. And for local machines it's not possible even when using a proxy server (website talks to a proxy, proxy talks to VNC server and tunnels the data to website).

---

#inspiration "Every day for years, Trollope reported in his “Autobiography,” he woke in darkness and wrote from 5:30 a.m. to 8:30 a.m., with his watch in front of him." (http://www.newyorker.com/magazine/2004/06/14/blocked). If I was able to do coding like this, I would have written many more programs than I have.

---

Digital footprints:
* http://discuss.bootstrapped.fm/t/advice-on-getting-free-users-to-pay-for-a-pro-option-of-a-product/5084/3?u=kjk
* https://www.indiehackers.com/forum/post/-Km7Y53UZUwjvKGZECh6?commentId=-Km7xYYmTwQJPBGQwFQo
* https://github.com/kurin/blazer/issues/17

2017-06-07

#go

Mostly finished local drive -> backblaze backup tool in Go. Backblaze is much cheaper at storage and bandwidth than S3 or Google Storage.

It's content-addressable store so I never re-upload files.

High-level architecture of the backup tool:
* you provide a directory to backup and list of extensions to backup (it's not a generic backup tool but a way to backup documents)
* the program traverses the directory recursively
* if a file has the right extension, calculate its sha1
* name of the file in backblaze is based on sha1 so it's O(1) to check if it's already there, in which case we skip the upload

Since the program is meant to be run many times over the same directories, I needed few tricks to make it faster:
* to avoid re-calculating sha1 of the file across the runs, I store previously calculated sha1 in a .csv "database". It's small enough that I can load existing database in memory on startup and update it by appending to .csv file as I calculate new sha1
* to speed up "is the file already there" check, I list all files in backblaze on startup, which gives me info about 1000 files per HTTP request (compared to more naive way of checking each file individually, which would require HTTP call for each file)

Go library I used for backblaze API had small bugs and inefficiencies, but the author fixed them (https://github.com/kurin/blazer/issues/16, https://github.com/kurin/blazer/issues/13, https://github.com/kurin/blazer/issues/14) within hours of reporting. Kudos because it's rare.

Also started on a very simple web front-end for listing the files. https://vuejs.org is great as a simple template library (have data -> render some html based on that). I know that Vue.js is much more than that, but it's perfect for simple templating as well. In the past I would laboriously create html as text in JavaScript or use some simple templating like mustache.

---

Playing a bit with https://sciter.com (trying to build a UI version of https://github.com/jonas/tig). Sciter is very promising in the sense of allowing to build UI quickly. However, it's badly documented. Good thing it has lots of examples in the SDK. #sciter
