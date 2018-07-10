Id: 16
Title: Extracting files from .7z archives in Go
Date: 2015-06-03T17:45:57-07:00
Format: Markdown
Tags: go, programming
Collection: go-cookbook
--------------

For a project written in Go I needed to access files inside .7z archive.

Unfortunately at this time there is no pure Go implementation for that.

An easy way to achieve this is to sub-launch `7z` executable to extract the file to disk and then proceed as usual.

I spent a little more effort and created a way to get an `io.ReadCloser` stream for a file inside .7z archive.

It still uses `7z` executable for the functionality but it doesn't create temporary file. Instead it streams decompressed data so it's more efficient.

I've packaged this as [lzmadec](https://github.com/kjk/lzmadec) library and this article describes how to use it.

Pre-requisites:

* install the library: `go get -u github.com/kjk/lzmadec`
* install `7z` binary:
	* on mac: `brew install p7zip`
	* on ubuntu: `apt-get install p7zip`
	* on windows: http://www.7-zip.org/

The library offers essentially 2 functions:

* get a list of all files and their metadata inside .7z archive
* get `io.ReadCloser` for a given file
* and as a convenience a function to extract a file to disk

Here's a skeleton of the code, error handling skipped for  clarity:

```go
var archive *lzmadec.Archive
archive, _ := lzmadec.NewArchive("foo.7z")

// list all files inside archive
for _, e := range archive.Entries {
	fmt.Printf("name: %s, size: %d\n", e.Path, e.Size)
}
firstFile := archive.Entries[0].Path

// extract to a file
archive.ExtractToFile(firstFile + ".extracted", firstFile)

// decompress to in-memory buffer
r, _ := archive.GetFileReader(firstFile)
var buf bytes.Buffer
_, _ = io.Copy(&buf, r)
// if not fully read, calling Close() ensures that sub-launched 7z executable
// is terminated
r.Close()
fmt.Printf("size of file %s after decompression: %d\n", firstFile, len(buf.Bytes()))
```

To see a more complete example: https://github.com/kjk/lzmadec/blob/master/cmd/test/main.go

If you're a Go programmer, you can use it as an example of using `os.Cmd` to capture stdout of a program for progressive reading. `7z` has `-so` options which tells it to send decompressed data to stdout. This is how to capture that as an `io.Reader`: https://github.com/kjk/lzmadec/blob/master/lzmadec.go#L225

Hopefully some day there will be pure Go library for accessing .7z files and this hack will obsolete. https://github.com/uli-go/xz is quite advanced but not there yet.
