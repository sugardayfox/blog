Id: 1407031
Title: Speeding up Go with custom allocators
Tags: go,programming
Date: 2012-11-26T12:17:53-08:00
Format: Markdown
HeaderImage: gfx/headers/header-07.jpg
Collection: go-cookbook
Description: Speeding up a well-known benchmark (construction of binary trees) ~4x by optimizing memory allocations.
---

**Summary**: using a custom allocator I was able to speed up an
allocation heavy program ([binary-trees
benchmark](http://benchmarksgame.alioth.debian.org/u64q/binarytrees.html))
~4x.

Allocation is expensive. It holds true for all languages. At the time of
this writing, Go (version 1.0.3) doesn’t have a garbage collector that
is as sophisticated as, say, garbage collector in JVM. There’s work
being done to improve Go’s GC but today allocations in Go are not as
cheap as they could be.

This can be seen in binary-trees benchmark which has almost no
computation but millions of allocations of small objects. As a result,
Java implementation is about 7x faster.

I was able to speed up Go code by about 4x by using a custom allocator.

The benchmark builds a large binary tree composed of nodes:

```go
type struct Node {
 int item
 left, right *Node
}
```

To allocate a new node we use `&Node{item, left, right}`.

Improving allocation-heavy code
-------------------------------

First, a correction. When I said that allocation is expensive, I
over-simplified.

In garbage-collected languages allocation is actually very cheap. In a
good allocator it’s just a single arithmetic operation to bump a
pointer, which is orders of magnitude cheaper than even the best
`malloc()` implementation. `Malloc()` has to maintain data structures to
keep track of allocated memory so that `free()` can return it back to
the OS.

More complicated reality is that it’s garbage collection phase that is
expensive.

Garbage collection (gc) is triggered every once in a while. It
recursively scans all the allocated objects, starting from known roots
and chasing pointers. It figures out which objects are not referenced by
any other object and frees them (this is the “garbage” in garbage
collection that has just been collected).

There are 2 insights we get from knowing how garbage collection works
internally:

-   the more objects there are, the more expensive garbage collection is
-   the more pointers we need to chase, the more expensive gc is

Knowing what the problem is, we know what the solution should be:

-   allocate less objects (e.g. by allocating them in bulk or re-using
    previously allocated objects)
-   don’t use pointers

As it happens, the majority of the 4x speedup I got in this particular
benchmark came from not using pointers

Speeding binary-trees shootout benchmark
----------------------------------------

We said that we should avoid pointers, so that garbage collector doesn’t
have to chase them. The new definition of `Node` struct is:

```go
type NodeId int

type struct Node {
 int item
 left, right NodeId
}
```

We changed `left` and `right` fields from `*Node` to an alias type
`NodeId`, which is just a unique integer representing a node.

What `NodeId` means is up to us to define and we define it thusly: it’s
an index into a `[]Node` array. That array is the backing store (i.e.
allocator) for our nodes.

When we need to allocate another node, we expand the `[]Node` array and
return the index to that node. We can trivially map `NodeId` to `*Node`
by doing `&array[nodeId]`.

Our implementation is a bit more sophisticated. In Go it’s easy to
extend the array with `append()` but it involves memory copy. We avoid
that by pre-allocating nodes in buckets and using an array of arrays for
storage. The code is still relatively simple:

```go
const nodes_per_bucket = 1024 * 1024

var (
  all_nodes [][]Node = make([][]Node, 0)
  nodes_left int = 0
  curr_node_id int = 0
)

func NodeFromId(id NodeId) *Node {
  n := int(id) - 1
  bucket := n / nodes_per_bucket
  el := n % nodes_per_bucket
  return &all_nodes[bucket][el]
}

func allocNode(item int, left, right NodeId) NodeId {
  if 0 == nodes_left {
    new_nodes := make([]Node, nodes_per_bucket, nodes_per_bucket)
    all_nodes = append(all_nodes, new_nodes)
    nodes_left = nodes_per_bucket
  }
  nodes_left -= 1
  node := NodeFromId(NodeId(curr_node_id + 1))
  node.item = item
  node.left = left
  node.right = right

  curr_node_id += 1
  return NodeId(curr_node_id)
}
```

Remaining changes to the code involve adding `NodeFromId()` call in a
few places.

You can compare
[original](https://github.com/kjk/kjkpub/blob/master/gobench/bintree.go)
to my [faster
version](https://github.com/kjk/kjkpub/blob/master/gobench/bintree3.go).

Another minor advantage if using integers instead of pointers in Node
struct is that on 64-bit machines we use only 4 bytes for an int vs. 8
bytes for a pointer.

Drawbacks of custom allocators
------------------------------

The biggest drawback is that we lost ability to free objects. Memory
we’ve allocated will never be returned to the OS until the program
exits.

It’s not a problem in this case, since the tree only grows and the
program ends when it’s done.

In different code this could be a much bigger issue. There are ways to
free memory even with custom allocators but they require more complexity
and evolve towards implementing a custom garbage collector at which
point it might be better to go back to simple code and leave the work to
built-in garbage collector.

How come Java is so much faster?
--------------------------------

It’s a valid question: both Java and Go have garbage collectors, why is
Java’s so much better on this benchmark?

I can only speculate.

Java, unlike Go, uses generational garbage collector, which has 2
arenas: one for young (newly allocated) objects (called nursery) and one
for old objects.

It has been observed that most objects die young. Generational garbage
collector allocates objects in nursery. Most collections only collect
objects in nursery. Objects that survived collections in nursery are
moved to the second arena for old objects, which is collected at a much
lower rate.

Go collector is a simpler mark-and-sweep collector which has to scan all
allocated objects during every collection.

Generational garbage collectors have overhead because they have to copy
objects in memory and update references between objects when that
happens. On the other hand they can also compact memory, improving
caching and they scan a smaller number of objects during each
collection.

In this particular benchmark there are many Node objects and they never
die, so they are promoted to rarely collected arena for old objects and
each collection is cheaper because it only looks at a small number of
recently allocated Node objects.

There’s hope for Go, though. The implementors are aware that garbage
collector is not as good as it could be and there is an ongoing work on
implementing a better one.

A win in C++ as well
--------------------

Optimizing by reducing the amount of allocations or making allocations
faster is applicable to non-gc languages as well, like C and C++,
because `malloc()\free()` are relatively slow functions.

Back in the day when I was working on Poppler, I achieved a significant
~19% speedup by [improving a string
class](/article/Performance-optimization-story.html)
to avoid an additional allocation in 90% of the cases. I now use this
trick in my C++ code e.g. in [SumatraPDF
code](https://code.google.com/p/sumatrapdf/source/browse/trunk/src/utils/Vec.h)

I also managed to improve Poppler by another \~25% by using a simple,
[custom allocator](https://bugs.freedesktop.org/show_bug.cgi?id=7910)

It's a good trick to know.
