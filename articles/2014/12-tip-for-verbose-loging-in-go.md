---
Id: 12
Title: Tip for verbose logging in Go
Date: 2014-12-03T18:21:29-08:00
Format: Markdown
Tags: programming,go
---

Go lacks a good ad-hoc debugging tool (i.e. a competent debugger).

For that reason we have to revert to caveman-style debugging: sprinkling the code with `fmt.Printf()`.

The problem with this approach is lack of selectivity: imagine you have 100 tests
and only 1 test fails. For debugging the issue you only need to see logs when
executing that 1 test but you'll drown in log output from all 100 tests.

My solution: control logging state with global variable `verboseLog` and allow toggling this flag per test.

Something like this:

```go
var (
    verboseLog = false
)

func myCodeWithBugs(s string) {
    if verboseLog {
        fmt.Printf("s: %s\n", s)
    }
    ...
}

func TestMyCode(t *testing.T) {
    var tests = []struct {
        ... test fields
        debug bool
    }{
        { ..., false }, // without verbose logging
        { ..., true },  // with verbose logging
    }
    for _, test := range tests {
        verboseLog = test.debug
        ...
    }
}
```
