Id: 18
Title: Optimizing JavaScript by using arrays instead of objects
Date: 2016-10-23T00:43:56-07:00
Tags: programming, javascript
Format: Markdown
--------------
@header-image gfx/headers/header-05.jpg

Best optimizations are achievied by thinking about a problem holistically.

In this article I describe an optimization that uses arrays instead of classes while providing a class API for accessing data.

Imagine you're building a web-based [note taking application](https://quicknotes.io).

It uses modern, single-page architecture. Front-end is written in React and backend provides JSON data to React.

The main view is a list of all notes of a given user. You need a backend api
call that returns list of user's notes. You survey how everyone else is
implementing such API and you come up with the following:
`/api/getnotes?user_id=<user_id>` call which returns JSON response that looks like:

```json
{
  "notes": [
    {
      "id": 1,
      "title": "first note",
      "createdAt": "2016-08-14 15:34:32Z",
      // ... more properties
    },
    {
      "id": 2,
      "title": "second note",
      "createdAt": "2016-08-14 16:03:12Z"
      // ... more properties
    },
    // ... more notes
  ]
}
```

You notice there's a lot of redundancy as we repeat property names in every note object.

In our case the structure of the note is fixed i.e. it always has the same properties. We can encode this data more efficiently:

```json
{
  "notes": [
    [1, "first note", "2016-08-14 15:34:32Z", ... more properties],
    [2, "second note", "2016-08-14 16:03:12Z", ... more properties],
    // ... more notes
  ]
}
```

This is a holistic optimization that achieves several speedups at once:

* backend generates less text (JSON response)
* backend compresses less data
* browser downloads less data
* browser decompresses less data
* browser parsers less JSON text
* JavaScript arrays are most likely more memory efficient that objects, so the
data uses less memory

If you know how gzip compression works you might protest that our effort to
remove property names is mostly futile because gzip is very good at removing
such redundancies.

I benchmarked [quicknotes.io](http://quicknotes.io) using my
own notes and found that even after compression the size difference of two
versions is ~50%. This might be a difference between a browser
downloading 150 kB of data vs. 300 kB.

This representation comes at a cost of programmer's convenience.

With objects we say `note.title`. With array representation it's more work:

```javascript
const noteIdIdx = 0;
const noteTitleIdx = 1;
// ... constants for more properties

const title = note[noteTitleIdx];
```

This is not great. We can improve this by writing accessor functions:

```javascript
function getTitle(note) {
  return note[noteTitleIdx];
}
```

JavaScript is a pliable language and we can get the best of both
worlds: array representation with class API.

We create a class that derives from `Array` and extends it with accessor functions.

 This example uses TypeScript, because static typing rocks, but will also work in pure JavaScript.

```javascript
class Note extends Array {
  constructor() {
    super();
  }

  ID(): int {
    return this[noteIdIdx];
  }

  Title(): string {
    return this[noteTitleIdx];
  }
  // .. more accessor functions
}

// this "upgrades" rawArray object from being Array instance
// to Note instance by patching prototype chain.
// Beware: if rawAray is not Array instance, bad things will happen
function toNote(rawArray: any): Note {
  Object.setPrototypeOf(rawArray, Note.prototype);
  return rawArray as Note;
}

// one way to convert raw array to object
const rawNote = [1, "first note", ... more properties];
const note = toNote(rawNote);
const title = note.Title();

// a less efficient but also less hacky way is by constructing a new Note object
const note = new Note(rawNote);
const title = note.Title();
```

Class `Note` extends built-in JavaScript `Array` so it's as efficient as an array and inherits all its functionality.

We add a couple of functions for getting/setting note data for a better API.

The magic happens in `toNote` function.

`rawNote` is an instance of `Array` We could add our accessor functions directly to `Array.prototype` but that would make them available to all `Array` instances.

By defining a class `Note` that inherits from `Array`, `Note.prototype` inherits all of `Array.prototype` functions and gets our additional functions.

In order to convert a raw array to `Note` object we can either construct a new object from raw array or "upgrade" the object with `Object.setPrototypeOf(note, Note.prototype)`.

This is dangerous: if the object being upgraded is not an instance of `Array`,
bad things will happen. It's not a technique that should be overused.

Upgrading the object in place should be more efficient than creating a new
object as it avoids an allocation.

On the other hand, [according to MDN](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/setPrototypeOf) changing prototype of an object makes for slowera access, so it can go either way.

To summarize:

* optimization is often achieved by looking at a problem as a whole
* thanks to flexibility of JavaScript we can implement a micro-optimization where we represent objects as arrays but add convenient class-like API
