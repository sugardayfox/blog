---
Id: 1Ll7
Title: Rotate log files daily in Go
Format: Markdown
Tags: for-blog, go
CreatedAt: 2017-06-14T02:07:56Z
UpdatedAt: 2017-07-09T20:35:12Z
HeaderImage: gfx/headers/header-17.jpg
Collection: go-cookbook
Description: How to rotate a log file daily.
Status: draft
---

If your program logs to a file, it's a good idea to rotate log files as that prevents them from growing too large.

After rotation you can backup them up to online storage or delete them.

Rotating daily is a good balance of how often to rotate vs. how large the log file can become.

You can write logs to stdout and use external program, like [logrotate](https://www.cyberciti.biz/faq/how-do-i-rotate-log-files/), to do the rotation.

I prefer the simplicity of handling that in my own code and Go's `io.Writer` interface makes it easy to implement a re-usable file ration functionality.

I did just that in package [dailyrotate](https://github.com/kjk/dailyrotate) ([documentation](https://godoc.org/github.com/kjk/dailyrotate)).

Here's basic use:

```go
import (
	"github.com/kjk/dailyrotate"
)

var (
	rotatedFile *dailyrotate.File
)

// called when file is closed. If didRotate is true, it was closed due to rotation
// at UTC midnight. Otherwise it was closed due to regular Close()
func onCloseHappened(path string, didRotate bool) {
	fmt.Printf("we just closed a file '%s', didRotate: %v\n", path, didRotate)
	if !didRotate {
		return
	}
	// we block writes until this returns so expensive processing should be done in
	// background goroutine
	go func() {
		// here you can implement things like:
		// - compressing rotated file
		// - deleting old files to free up disk space
		// - upload rotated file to backblaze/google storage/s3 for backup
		// - analyze the content of the file
	}()
}

func initRotatedFileMust() {
	pathFormat := filepath.Join("dir", "2006-01-02.log")
	w, err := dailyrotate.NewFile(pathFormat, onCloseHappened)
	panicIfErr(err)
}

func logString(s string) error {
	_, err = io.WriteString(rotatedFile, s)
	panicIfErr(err)
}

func shutdownLogging() {
	rotatedFile.Close()
}
```

Here's a [real-life example](https://github.com/kjk/blog/blob/ee30c22379c90642880c8fae33fa3b767a22cb64/visitor_analytics.go#L229) of processing rotated file:


`dailyrotate.File` is `io.Writer` and safe to use from multiple goroutines.

In addition to `Write(d []byte)` It also implements a `Write2(d []byte, flush bool) (string, int64, int, error)`. It has 2 improvements over `Write`:

* allows to flush in a single call. Flushing after each write is slower but less likely to loose data or corrupt the file when program crashes
* it returns file path and offset in the file at which the data was written. This is important for building a random-access index to records written to `dailyrotate.File`

## Other real-world uses

Rotation is not limited to log files. I use it as part of poor-man's [web server analytics system](https://github.com/kjk/blog/blob/master/visitor_analytics.go).

I log info about web requests to `dailyrotate.File` using my [siser](https://github.com/kjk/siser) simple serialization format.

When a file is rotated, I compress it, upload to backblaze for backup and delete local files older than 7 days to free up space.

I also calculate basic statistics for the day and e-mail a summary to myself.

I know, I could just use Google Analytics, and I do. My little system has advantages:
* it tells me about missing pages (404). That alerts me if break something. Also people sometimes link incorrectly to my website and knowing what are bad links, I can add re-directs for them
* it e-mails me daily the summary of most important information. I keep an eye on things with minimal effort.
