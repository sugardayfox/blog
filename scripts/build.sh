#!/bin/bash
set -u -e -o pipefail

go build -o blog_app *.go
rm blog_app
