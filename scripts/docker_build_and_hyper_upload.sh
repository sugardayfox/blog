#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

rm -rf $TMPDIR/gdep
GOOS=linux GOARCH=amd64 gdep go build -o blog_app_linux

docker build --tag kjksf/blog:latest --tag blog:latest .
docker push kjksf/blog:latest
hyper pull kjksf/blog:latest
