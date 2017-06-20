Id: JyRZ
Title: Generating good, random and unique ids in Go
Format: Markdown
Tags: for-blog, go, draft
CreatedAt: 2017-06-20T09:07:51Z
UpdatedAt: 2017-06-20T09:07:52Z
--------------
@header-image gfx/headers/header-01.jpg

Imagine you're writing a [note taking](https://quicknotes.io) application.

Each note must have a unique id.

Generating unique ids is easy if you can coordinate.

The simplest way is to get database to do it: use `AUTOINCREMENT` column and the database will generate unique id when you insert the note into a database.

What if you can't coordinate?

For example, you want the app to also generate unique note id when offline, when it cannot talk to the database.

The requirement of non-coordinated generation of unique ids comes up often in distributed systems.

A simple solution is to generate a random id. If you give it 16 bytes of randomness, the chances of generating the same random number are non-existent.

It's such a common problem that we have a standard for this called [UUID/GUID](https://en.wikipedia.org/wiki/Universally_unique_identifier), created over 30 years ago.

We can do better than GUID. A good random unique id:
* is random; we can't skip the basics
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
[github.com/oklog/ulid](https://github.com/oklog/ulid) | **`01BJMVNPBBZC3E36FJTGVF0C4S`** | 6 bytes time (milliseconds) + 8 bytes random

You can see how the values change over time by refreshing [this test page](/tools/generate-unique-id) a couple of times.

How to generate unique ids using different libraries:

```go
import (
	"github.com/kjk/betterguid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

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

func main() {
	genXid()
	genKsuid()
	genBetterGUID()
	genUlid()
}
```

Full example: [generate-unique-id/main.go](https://github.com/kjk/go-cookbook/blob/master/generate-unique-id/main.go)

Which one to use? All of them are good.

`oklog/ulid` allows custom entropy (randomness) source but pays for that with complex API.


To learn more:
* https://segment.com/blog/a-brief-history-of-the-uuid/ : the history of ksuid
* https://firebase.googleblog.com/2015/02/the-2120-ways-to-ensure-unique_68.html : inspiration for betterguid design
* http://antoniomo.com/blog/2017/05/21/unique-ids-in-golang-part-1/
* http://antoniomo.com/blog/2017/05/28/unique-ids-in-golang-part-2/
* http://antoniomo.com/blog/2017/06/03/unique-ids-in-golang-part-3/
