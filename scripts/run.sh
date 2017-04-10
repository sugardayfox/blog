#!/bin/bash
set -u -e -o pipefail

echo "building"
go build -o blog_app *.go
#go build -o blog_app *.go
echo "running in `pwd`"
./blog_app -addr=localhost:5020
rm blog_app
