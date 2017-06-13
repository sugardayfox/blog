#!/bin/bash
set -u -e -o pipefail

go build -o blog_app -ldflags "-X main.sha1ver=`git rev-parse HEAD`"
rm blog_app
