Id: 3
Title: Analyzing browserify bundles to minimize JavaScript bundle size
Date: 2017-01-04T22:02:27-08:00
Format: Markdown
--------------
When building web apps, it's important to keep the size of JavaScript code delivered to the browser as small as possible.

I write in ES6 or TypeScript then use browserify to combine all JavaScript code
into a single bundle file. For production builds I use uglify to make the bundle smaller.

Unfortunately, by default we are blind to what ends up in the final bundle.
A single `import` can introduce surprising, unneeded dependencies.

First step of fixing bloat is to see what code ends up in the final bundle.

## Disc

[Disc](http://hughsk.io/disc/) is one tool that visualizes the content of JavaScript bundle.

To use it:

* `npm install -g disc`
* add `fullPaths: true` option to `browserify` plugin (without it file paths are turned into opaque numbers)
* `discify dist/bundle.min.js >out.html` (or whatever `bundle.min.js` is called in your setup)
* `open out.html` (on mac, or open manually in the browser)

The visualization is very pretty but not very good for understanding.

## source-map-explorer

[source-map-explorer](https://www.npmjs.com/package/source-map-explorer) shows the same information but in a more useful way.

To use it:

* `npm install -g source-map-explorer`
* make sure that you generate JavaScript maps file
* `source-map-explorer dist/bundle.min.js dist/bundle.min.js.map`

This will open the browser for you with the treemap visualization.

## Analyzing dependency tree

Disc and source-map-explorer can tell you what but not why.

When you see a JavaScript package that shouldn't be there, you need to know why
it's there i.e. where it was imported from.

I haven't found a tool that makes it easy, but it's possible to create
a primitive debug tool yourself.

```javascript
var through = require('through2'),

var b = browserify(browserifyOpts);
if (showDeps) {
  // for debugging dump (flattened and inverted) dependency tree
  // b is browserify instance
  b.pipeline.get('deps').push(through.obj(
  function(row, enc, next) {
    // format of row is { id, file, source, entry, deps }
    // deps is {} where key is module name and value is file it comes from
    console.log(row.file || row.id);
    for (let k in row.deps) {
      const v = row.deps[k];
      console.log('  ', k, ':', v);
    }
    next();
  }));
}
```

This displays dependencies in the format:
```
/quicknotes/node_modules/react-dom/lib/LinkedValueUtils.js
   ./reactProdInvariant : /quicknotes/node_modules/react-dom/lib/reactProdInvariant.js
   ./ReactPropTypesSecret : /quicknotes/node_modules/react-dom/lib/ReactPropTypesSecret.js
   react/lib/React : /quicknotes/node_modules/react/lib/React.js
   fbjs/lib/invariant : /quicknotes/node_modules/fbjs/lib/invariant.js
   fbjs/lib/warning : /quicknotes/node_modules/fbjs/lib/warning.js
```

It's not an ideal presentation but you can figure out who ultimately imports
a given JavaScript file by chasing chain of imports.

## Things I learned

How does it help in practice? Here are 2 examples of how I reduced JavaScript bundle bloat by using those tools.

### bloated highlight.js

In [QuickNotes](https://quicknotes.io) I use [highlight.js](https://highlightjs.org/) library to do syntax highlighting for code snippets.

Looking at output of source-map-explorer I noticed that highlight.js is
476 kB in size. That seemed excessive.

The problem was that while core of highlight.js is small, it supports 168 languages
and doing `import 'highlight.js'` would bundle all of.

I only need to support small subset of most popular languages.

One way to fix it would be to use https://highlightjs.org/download/ to generate
a custom bundle. That would require repeating this manual step when I want
to use the newer version.

I settled on a hacky but more automated solution.

Doing `import 'highlight.js'` loads `node_modules/highlight.js/index.js` which
imports all languages.

I created a custom `index.js` that only imports the languages I want. Bbefore
every compilation, I over-write `node_modules/highlight.js/index.js` with my
custom version.

That way I can still use npm to manage the library and easily update to new version.

The result? Saved 416 kB.

### bloated seedrandom.js

At [work](https://wwww.folsomlabs.com) we use tiny [seendrandom.js](https://github.com/davidbau/seedrandom/blob/released/seedrandom.js) library.

When inspecting our JavaScript bundle I noticed suspicious libraries in it,
like asn1 decoder.

I suspected our code doesn't do asn1 decoding. Searching the codebase didn't
turn up any direct use.

I speculated that it's imported indirectly by some other library.

I used my ad-hoc dependency tree dump to figure out that this code is
imported from `seedrandom.js`.

This piece of code is a culprit:

```javascript
  // When in node.js, try using crypto package for autoseeding.
  try {
    nodecrypto = require('crypto');
  } catch (ex) {}
```

Since node libraries are available during build step this line adds 294 kB of
unneeded crypto code to our web app.

The fix was to fork the repo and remove those lines.

## Automating things

It's handy to be able to re-run this analysis. Here's a sample script
`analyze_bundle.sh` I have in one of my projects:
```bash
#!/usr/bin/env bash
set -u -e -o pipefail

# uses source-map-explorer (https://www.npmjs.com/package/source-map-explorer)
# to visualize what modules end up in final javascript bundle.

install_sme() {
  if [ ! -f ./node_modules/.bin/source-map-explorer ]; then
    npm install source-map-explorer
  fi
}

analyze_prod()
{
  rm -rf s/dist/*
  install_sme

  ./node_modules/.bin/gulp jsprod
  ./node_modules/.bin/source-map-explorer s/dist/bundle.min.js s/dist/bundle.min.js.map
}

analyze_prod
```

The particulars will depend on your build system. The general idea is to run
the build to generate `.js` and `.map.js` files and run `source-map-explorer`
for analysis.
