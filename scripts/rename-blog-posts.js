#!/usr/bin/env node

const child_process = require('child_process');
const fs = require('fs');
const walk = require('walk');

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

walker.on('end', rename_files);

// YYYY-DD/foo.txt -> YYYY/DD-foo.txt
function path_rename(path) {
  let parts = path.split("/");
  const yyyy_dd = parts[0];
  const name = parts[1];
  parts = yyyy_dd.split("-");
  if (parts.length !== 2) {
    return null;
  }
  const yyyy = parts[0];
  const dd = parts[1];
  return `${yyyy}/${dd}-${name}`;
}

function dir_for_path(path) {
  const parts = path.split("/");
  return parts[0];
}

function rename_git_files(to_rename) {
  const paths = to_rename.pop();
  if (!paths) {
    return;
  }
  const [path, new_path] = paths;
  const cmd = `git mv ${path} ${new_path}`;
  console.log(`cmd: ${cmd}`);
  child_process.exec(cmd, (error, stdout, stderr) => {
    if (error) {
      console.log(`error: ${error}, stdout: ${stdout}, stderr: ${stderr}`);
      return;
    }
    rename_git_files(to_rename);
  });
}

function mkdirIfNotExistsSync(dir) {
    if (!fs.existsSync(dir)) {
      console.log(`created ${dir}`);
      fs.mkdirSync(dir);
    }
}

function rename_files() {
  files.sort();
  const dirsDict = {};
  const to_rename = [];
  files.forEach(path => {
    const new_path = path_rename(path);
    if (new_path) {
      const dir = dir_for_path(new_path);
      dirsDict[dir] = 1;
      to_rename.push([path, new_path]);
    }
  });

  const dirs = Object.keys(dirsDict);
  dirs.forEach(mkdirIfNotExistsSync);

  rename_git_files(to_rename);
}


