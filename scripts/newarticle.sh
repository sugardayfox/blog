#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

if [ -e blog_app ]
then
    ./blog_app -newarticle="$*"
else
    go build -o blog_app *.go
    ./blog_app -newarticle="$*" || true
    rm blog_app
fi
