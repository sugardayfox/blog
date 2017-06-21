Id: JyRZ
Title: Generating good, random and unique ids in Go
Format: Markdown
Tags: for-blog, go, published
CreatedAt: 2017-06-21T08:51:55Z
UpdatedAt: 2017-06-21T08:51:56Z
--------------
@header-image gfx/headers/header-01.jpg

Imagine you're writing a [note taking](http://quicknotes.io) application.

Each note needs a unique id.

Generating unique ids is easy if you can coordinate.

The simplest way is to get database to do it: use `AUTOINCREMENT` column and the database will generate unique id when you insert a new `note` row into a table.

What if you can't coordinate?

For example, you want the app to also generate unique note id when offline, when it cannot talk to the database.

The requirement of non-coordinated generation of unique ids comes up often in distributed systems.

A simple solution is to generate a random id. If you give it 16 bytes of randomness, the chances of generating the same random number are non-existent.

It's such a common problem that over 30 years ago we created a standard for this called [UUID/GUID](https://en.wikipedia.org/wiki/Universally_unique_identifier).

We can do better than GUID. A good random unique id:
* is unique; we can't skip the basics
* can be sorted by its string representation
* is time-clustered i.e. ids generated at the same time are close to each other when sorted
* string representation can be used as part of URL without escaping
* the shorter, the better

There are few Go implementation of such id, following the same basic idea:
* use time as part of the id to achieve time-clustering
* fill rest of the id with random data
* encode as a string in a way that allows lexicographic sorting and is url-safe

Here are Go packages for generating unique id and how their ids look like in string format:


package | id | format
--- | --- | ---
[github.com/segmentio/ksuid](https://github.com/segmentio/ksuid) | **`0pPKHjWprnVxGH7dEsAoXX2YQvU`** | 4 bytes of time (seconds) + 16 random bytes
[github.com/rs/xid](https://github.com/rs/xid) | **`b50vl5e54p1000fo3gh0`** | 4 bytes of time (seconds) + 3 byte machine id + 2 byte process id + 3 bytes random
[github.com/kjk/betterguid](https://github.com/kjk/betterguid) | **`-Kmdih_fs4ZZccpx2Hl1`** | 8 bytes of time (milliseconds) + 12 random bytes
|[github.com/sony/sonyflake](https://github.com/sony/sonyflake)| **`20f8707d6000108`**| ~6 bytes of time (10 ms) + 1 byte sequence + 2 bytes machine id|
[github.com/oklog/ulid](https://github.com/oklog/ulid) | **`01BJMVNPBBZC3E36FJTGVF0C4S`** | 6 bytes of time (milliseconds) + 8 bytes random
|[github.com/satori/go.uuid](https://github.com/satori/go.uuid)| **`5b52d72c-82b3-4f8e-beb5-437a974842c`** | UUIDv4 from [RFC 4112](http://tools.ietf.org/html/rfc4122) for comparison

You can see how the values change over time by refreshing [this test page](/tools/generate-unique-id) a couple of times.

How to generate unique ids using different libraries:

```go
import (
	"github.com/kjk/betterguid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"
)

// To run:
// go run main.go

func genXid() {
	id := xid.New()
	fmt.Printf("github.com/rs/xid:           %s\n", id.String())
}

func genKsuid() {
	id := ksuid.New()
	fmt.Printf("github.com/segmentio/ksuid:  %s\n", id.String())
}

func genBetterGUID() {
	id := betterguid.New()
	fmt.Printf("github.com/kjk/betterguid:   %s\n", id)
}

func genUlid() {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	fmt.Printf("github.com/oklog/ulid:       %s\n", id.String())
}

func genSonyflake() {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	// Note: this is base16, could shorten by encoding as base62 string
	fmt.Printf("github.com/sony/sonyflake:   %x\n", id)
}

func genUUIDv4() {
	id := uuid.NewV4()
	fmt.Printf("github.com/satori/go.uuid:   %s\n", id)
}

func main() {
	genXid()
	genKsuid()
	genBetterGUID()
	genUlid()
	genSonyflake()
	genUUIDv4()
}
```

Full example: [generate-unique-id/main.go](https://github.com/kjk/go-cookbook/blob/master/generate-unique-id/main.go)

Which one to use? All of them are good.

I would pick either `rs/xid` or `segmentio/ksuid`.

`oklog/ulid` allows custom entropy (randomness) source but pays for that with complex API.

`sony/sonyflake` is the smallest but also the least random. It's based on Twitter's design for generating IDs for tweets.

For simplicity the example code serializes `sony/snoflake` in base16. It would be even shorter in base62 encoding used by other libraries, but other libraries provide that out-of-the-box and for `sony/snoflake` I would have to implement it myself.

The last one is UUID v4 from [RFC 4112](http://tools.ietf.org/html/rfc4122), for comparison.

To learn more:
* https://segment.com/blog/a-brief-history-of-the-uuid/ : the history of ksuid
* https://firebase.googleblog.com/2015/02/the-2120-ways-to-ensure-unique_68.html : inspiration for betterguid design
* http://antoniomo.com/blog/2017/05/21/unique-ids-in-golang-part-1/
* http://antoniomo.com/blog/2017/05/28/unique-ids-in-golang-part-2/
* http://antoniomo.com/blog/2017/06/03/unique-ids-in-golang-part-3/
* https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake.html : details of Twitter's snowflake on which sonyflake is based
