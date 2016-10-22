#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

GOOS=linux GOARCH=amd64 go build -o blog_app_linux

docker build --tag kjksf/blog:latest --tag blog:latest .
