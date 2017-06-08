2017-06-08

1. Started this worklog. More code needed to actually show it.

2. Reading more blog posts from http://thestartuptoolkit.com/.

Insight: if you're business is "better X", use X to make sure you're really improving on it (http://thestartuptoolkit.com/blog/2011/09/use-the-tools-youre-displacing/) #business

3. #read https://blog.figma.com/webassembly-cut-figmas-load-time-by-3x-76f3f2395164. Summary: webassembly is much faster than asm.js. #webassembly

4. #idea For faster PDF rendering, re-compile it. PDF format is rather complicated, mostly text format. A PDF renderer first has to parse it into in-memory representation and then renders it.

Imagine a binary format that is designed for fast loading and one-time PDF => binary format compilation step. With the right format, parsing step could essentially be eliminated.

We could even improve rendering time with some additional optimizations.

For example, there are PDF documents with gigantic images (think 4000 x 4000 pixels) that could be resized to some reasonable value and saved in a format that is fastest to decode.

I'm sure there are optimizations possible for complex vector graphics (e.g. removing stuff that we know will be invisible).

5. Digital footprints:
* http://discuss.bootstrapped.fm/t/advice-on-getting-free-users-to-pay-for-a-pro-option-of-a-product/5084/3?u=kjk
* https://www.indiehackers.com/forum/post/-Km7Y53UZUwjvKGZECh6?commentId=-Km7xYYmTwQJPBGQwFQo

2017-06-07

1. Mostly finished local drive -> backblaze backup tool in Go. Backblaze is much cheaper at storage and bandwidth than S3 or Google Storage.

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

2. Playing a bit with https://sciter.com (trying to build a UI version of https://github.com/jonas/tig). Sciter is very promising in the sense of allowing to build UI quickly. However, it's badly documented. Good thing it has lots of examples in the SDK.
