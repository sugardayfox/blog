Id: 18
Title: JavaScript optimization by using arrays instead of objects
Date: 2016-10-23T00:43:56-07:00
Tags: programming, javscript
Format: Markdown
--------------
This article talks about two things:

- best optimizations are achievied by thinking about a problem holistically
- In JavaScript, a clever micro-optimization that uses arrays instead of classes while preserving a class-based interface for accessing data

Imagine you're building a note-taking web application. It uses modern, single-page architecture with front-end written in React and backend that mostly implements API calls returning JSON data that drive the display.

The main view is a list of all notes of a given user. You need a backend api call that returns list of user's notes. You survey how everyone else is implementing such API and you come up with the following: `/api/getnotes?user_id=<user_id>` call returns JSON repronse that looks like:

```javascript
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

You notice there's a lot of redundancy as we repeat key names for every object. In our case the structure of a note is fixed i.e. at always has the same set of properties. We can encode this data much more efficiently:

```javascript
{
  "notes": [
  [1, "first note", "2016-08-14 15:34:32Z", ... more properties],
  [2, "second note", "2016-08-14 16:03:12Z", ... more properties],
  // ... more notes
  ]
}
```

This is a holistic optimization that creates several speedups at once:

- backend generates less text (JSON response)
- backend compresses less data
- browser downloads less data
- browser decompresses less data
- browser parsers less JSON text
- JavaScript arrays are most likely more memory efficient that objects, so the data uses less memory

If you know how gzip compression works you might protest that our effort to remove property names is mostly futile because gzip is very good at removing such redundancies. As with all benchmarks the exact results depend on exact details but when benchmarking [quicknotes.io](http://quicknotes.io) using my own notes, I found that even after compression the size difference of two compressed versions is ~50%. This might be a difference between a browser downloading 150kB of data vs. 300kB.

This does come at a small cost: accessing the data in JavaScript is less convenient. With objects we can say: `note.title` . With array representation, we have to do more:

```javascript
const noteIdIdx = 0;
const noteTitleIdx = 1;
// ... constants for more properties

const title = note[noteTitleIdx];
```

This is not great. We can improve that by writing accessor functions:

```javascript
function getTitle(note) {
  return note[noteTitleIdx];
}
```

Thankfully, JavaScript is very dynamic language so we can get the best of both worlds: array representation but class interface. This example uses TypeScript, because static typing rocks, but will work in pure ES2015.

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
// to Note instance. Note: if rawAray is not Array instance,
// bad things will happen
function toNote(rawArray: any): Note {
  rawArray.__proto__ = Note.prototype;
  return rawArray as Note;
}

// now we can can treat a raw array as Note object by patching class type
// of existing object
const rawNote = [1, "first note", ... more properties];
const note = toNote(rawNote);
const title = note.Title();

// a less efficient but also less hacky way is by constructing a new Note object
const note = new Note(rawNote);
const title = note.Title();
```

Class Note extends built-in JavaScript Array so it's as efficient as JavaScript array and inherits all its functionality. We add a couple of functions for getting/setting note data. That gives us efficiency and good interface.

The magic happens in `toNote` function. Since `rawNote` is an instance of `Array` , `rawNote.__proto__` is `Array.prototype` . We could add our accessor function directly to `Array.prototype` but that would make them available to all `Array` instances, even those that are not notes.

By definining a class `Note` that inherits from `Array` , `Note.prototype` inherits all of `Array.prototype` functions and gets our additional functions. In order to convert a raw array to `Note` object we can either construct a new object from raw array or "upgrade" the object by replacing `__proto__` property. This is dangerous: if the object being upgraded is not an instance of `Array` bad things will happen. It's not a technique that should be over-used.

Upgrading the object in place should be more efficient than creating a new object because it avoids an allocation. However [according to MDN](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/proto) , updating `__proto__` and subsequent operations on such objects are slow so it can go either way.

To summarize:

- optimization is often achieved by looking at a problem as a whole
- thanks to flexibility of JavaScript we can implement a micro-optimization where we represent objects as arrays but preserve convenient class-like interface to an object
