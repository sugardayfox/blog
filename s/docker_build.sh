#!/bin/bash
set -ue -o pipefail -o verbose

# fail if css has changed. This is called by deploy script so it'll stop deploy
go run ./s/css-update.go

dep ensure
GOOS=linux GOARCH=amd64 go build -o blog_linux -ldflags "-X main.sha1ver=`git rev-parse HEAD`"
docker build --no-cache --tag blog:latest .
