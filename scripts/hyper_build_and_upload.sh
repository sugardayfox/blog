#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

GOOS=linux GOARCH=amd64 go build -o blog_app_linux
echo "docker build"
docker build --tag blog:latest .
echo "docker save"
docker save blog:latest | gzip | aws s3 cp - s3://kjkpub/tmp/blog.tar.gz
echo "hyper load"
hyper load -i $(aws s3 presign s3://kjkpub/tmp/blog.tar.gz)
aws s3 rm s3://kjkpub/tmp/blog.tar.gz
