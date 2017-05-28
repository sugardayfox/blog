#!/bin/bash
set -ue -o pipefail -o verbose

GOOS=linux GOARCH=amd64 go build -o blog_linux -ldflags "-X main.sha1ver=`git rev-parse HEAD`"
docker build --no-cache --tag blog:latest .
