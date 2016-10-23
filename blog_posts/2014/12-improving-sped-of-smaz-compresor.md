Id: 13
Title: Improving speed of pure Go SMAZ compressor by 2.6x/1.5x
Date: 2014-12-11T00:19:49-08:00
Format: Markdown
Tags: go, programming
--------------

I was testing fast compressors in pure Go. One of them was [Go implementation](https://github.com/kjk/smaz)
of [SMAZ algorithm](https://github.com/antirez/smaz) for compressing small
strings. It's simple, fast and works well for English text.

It wasn't as fast as I expected so I looked at the code and with a few tweaks
I managed to speed up decompression 2.61x times and compression 1.54x times:

```
kjkmacpro:smaz kjk$ benchcmp before.txt after.txt
benchmark                  old ns/op     new ns/op     delta
BenchmarkCompression       3387936       2195304       -35.20%
BenchmarkDecompression     2667583       1022908       -61.65%

benchmark                  old MB/s     new MB/s     speedup
BenchmarkCompression       40.35        62.26        1.54x
BenchmarkDecompression     28.34        73.90        2.61x
```

The speed increase came from 3 micro-optimizations.

## 1. Don't use `bytes.Buffer` if `[]byte` will do.

The biggest decompression speed-up came from
[this change](https://github.com/kjk/smaz/commit/7adaf22db621f66027e38bd1ee4d36f351025043) where I replaced the use of `bytes.Buffer` with using slices directly.

`bytes.Buffer` is a wrapper around `[]byte`. It adds convenience by
implementing popular interfaces like `io.Reader`, `io.Writer` etc. but
decreased speed is the price of that.

Usuaully it doesn't matter but when there are lots of operations on `byte.Buffer`,
even small differences add up.

## 2. Re-using buffers is another common optimization trick in Go.

The original API was:
```
compressed := smaz.Compress(source)
```

`Compress` function has no option but to allocate a new buffer for the compressed
data every time. Allocations are not free and they slow down the program by
making garbage collector do more work.

Other compression libraries allow the caller to provide a buffer for the result:
```
compressed := make([]byte, 1024)
compressed = smaz.Encode(compressed, source)
```

If the buffer is not big enough, it'll be enlarged. If the caller
doesn't want additional complexity of managing reusable buffers, it can
pass `nil`.

## 3. Avoid un-necessary copies

Compression and decompression improves reading data from memory, transforming it
and writing the result to another memory location.

Memory operations are expensive. You can execute [7 CPU instructions](https://gist.github.com/kjk/0cd9e13e8b5f1046b697) for one memory operation in L2 cache.

I noticed that compression was making unnecessary temporary copies of data.
The code got [a bit more complicated](https://github.com/kjk/smaz/commit/754db648b7cd39fb12120a851e3d1106d2dff3e0) but also 1.14x faster.

## A digression on benchmarking tools in Go

One of the features that distinguish Go from other programming language
implemenations is that out-of-the-box it comes with tooling for testing,
profiling and benchmarking.

Writing benchmarks is straightforward. Here's a benchmark for compression speed:

```go
func BenchmarkCompression(b *testing.B) {
    b.StopTimer()
    inputs, n := loadTestData(b)
    b.SetBytes(n)
    b.StartTimer()
    var dst []byte
    for i := 0; i < b.N; i++ {
        for _, input := range inputs {
            dst = Encode(dst, input)
        }
    }
}
```

You run the benchmarks with `go test -bench=.`. You can benchmark only selected
function thanks to `-bench` argument (or pass `.` to run all of them).

Go minimizes amount of work the programmer needs to do in several ways:

* benchmarking functions are automatically recognized by convention: a function
that starts with `Benchmark` in `*_test.go` file is a benchmark function
* the results are in a standardized, human-readable form
* benchmarking tool not only measures time but you can also get MB/s metric
by using `testing.B.SetBytes()`, which is a better way to think about and compare
code like compression algorithms.

Finally, Go comes with a tool that makes it easy to compare speed before and
after the change:
```
> go get -u golang.org/x/tools/cmd/benchcmp
> go test -bench=. >before.txt
> ... make the changes
> go test -bench=. >after.txt
> benchcmp before.txt after.txt
```
