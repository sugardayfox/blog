Id: vkeR
Title: Simple serialization format for logging and analytics in Go
Format: Markdown
Tags: for-blog, draft, go
CreatedAt: 2017-07-09T03:59:28Z
UpdatedAt: 2017-07-09T20:46:20Z
--------------
@header-image gfx/headers/header-14.jpg
@collection go-cookbook
@description Usage, design and implementation of simple serialization format for logging and analytics in Go.
@status draft

I was looking to save multiple records with somewhat flexible schema to a file. Go has plenty of options but I couldn't find anything that was just right.

My desired features were:
* human-readable so that I can use tools like `grep` or `tail` to look at the file
* records are not fixed i.e. can have variable number of fields
* no need to support nested records
* allows for simple and efficient implementation

Here are some most popular available options and why they don't exactly fit the bill:
* csv : efficient but uses fixed records and not very readable if there are many fields on a single line
* json : not very readable
* protocol buffers : binary so not human readable, needs up-front scheme

I designed and implemented my own format in [siser](https://github.com/kjk/siser) library. 

You'll not be surprised by how it looks:

```
url: http://blog.kowalczyk.info/index.html
code: 200
large field:+13789
this is large data, 13789 bytes in size...
---
```

It's your standard key/value serialization with one neat feature: support for large data (e.g. an image or long text).

It's line oriented format, where each line is `${key}: ${value}\n`.

If the value is larger than 120 bytes or is not ascii text (bytes are outside of 32-127 range), we serialize it as large value:
```
${key}:+${value_length}\n
${value}\n
```

To separate records, we use `---\n`.

I use this format for 2 main purposes:
* structured logging to help in debugging
* logging events for analytics

## Using the library

When I use it for logging events for analytics, each record correspond to an event. I save the events to a file and later on process the whole file record-by-record and calculate desired statistics.

Here's how we would log info about HTTP requests:
```go
func logHTTPRequest(w io.Writer, url string, ipAddr string, statusCode int) error {
	var r siser.Record
	// you can append multiple key/value pairs at once
	r.Append("url", url, "ipaddr", ipAddr)
	// or assemble with multiple calls
	r.Append("code", strconv.Itoa(statusCode))
	d := r.Marshal()
	_, err := w.Write(d)
	return err
}
```

Here's a [full example](https://github.com/kjk/blog/blob/b18317d3dbde1d21745aaea615d952f2c2e158c8/visitor_analytics.go#L309) of logging HTTP requests.

Let's say we wrote the data to `http_access.log` file. Here's how we would process the records:
```go
f, err := os.Open("http_access.log")
panicIfErr(err)
defer f.Close()
r := siser.NewReader(f)
for r.ReadNext() {
	record := r.Record()
	code, ok := r.Get("code")
	// get rest of values and do something with them
}
panicIfErr(r.Err())
```

Here's a [full example](https://github.com/kjk/blog/blob/b18317d3dbde1d21745aaea615d952f2c2e158c8/visitor_analytics.go#L108) of calculating basic daily statistics on HTTP requests (most frequently visited pages, most frequent 404s, most frequent referers).

## Annoyances

Simplicity and speed has a cost.

The library doesn't offer marshaling directly to/from structs. I could add it but reflection is slow so I don't want to encourage its use.

It also doesn't directly support arbitrary types like `int` or `time.Time` so you have to convert them to/from string yourself.

Again, I could support for it. Instead of using `[]struct` I could use `[]interface{}` and support most common primitive types, but boxing as `interface{}` is more expensive than not doing that, so again I avoid slow-by-default solution at the cost of slightly worse API.

## Implementation notes

Some implementation decisions were made with performance in mind.

Given key/value nature of the record, an easy choice would be to use `map[string]string` as source to encode/decode functions.

However `[]string` is more efficient than a `map`. Additionally, a slice can be reused across multiple records. We can clear it by setting the size to zero and reuse the underlying array. A map would require allocating a new instance for each record, which would create a lot of work for garbage collector.

When serializing, you need to use `Reset` method to get the benefit of efficient re-use of the record.

When reading and deserializing records, `siser.Reader` uses this optimization internally.

The format avoids the need for escaping keys and values, which helps in making encoding/decoding fast. 

How does that play out in real life? I wrote a [benchmark](https://github.com/kjk/siser/blob/293341408be76f2b40b3f64b2c78de61bb3a887e/serialize_test.go#L132) comparing siser vs. `json.Marshal`. It's about 30% faster:
```
$ go test -bench=.
BenchmarkSiserMarshal-8   	 1000000	      1903 ns/op
BenchmarkJSONMarshal-8    	  500000	      2905 ns/op
```

The format is binary-safe and works for serializing large values e.g. you can use png image as value.

It's also very easy to implement in any language.

Code for this chapter: https://github.com/kjk/go-cookbook/tree/master/simple-serialization
