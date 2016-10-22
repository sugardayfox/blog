#!/usr/bin/env node

const walk    = require('walk');
const files   = [];

function clean_path(path) {
  if (path.startsWith("./")) {
    return path.substr(2);
  }
  return path;
}

process.chdir("blog_posts");
const walker  = walk.walk(".", { followLinks: false });
walker.on('file', function(root, stat, next) {
  let path = clean_path(root + "/" + stat.name);
  files.push(path);
  next();
});

// YYYY-DD/foo.txt -> YYYY/DD-foo.txt
function path_rename(path) {
  let parts = path.split("/");
  const yyyy_dd = parts[0];
  const name = parts[1];
  parts = yyyy_dd.split("-");
  const yyyy = parts[0];
  const dd = parts[1];
  return `${yyyy}/${dd}-${name}`;
}

walker.on('end', function() {
  files.sort();
  for (const path of files) {
    const new_path = path_rename(path);
    console.log(`${path} => ${new_path}`);
    // TODO: run git rename
  }
});
