#!/bin/bash
set -u -e -o pipefail

GOOS=linux GOARCH=amd64 go build -o blog_app_linux

docker build --tag blog:latest .
docker run --rm -it -v ~/data/blog:/data -p 5020:80 blog:latest
