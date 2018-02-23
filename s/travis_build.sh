#!/bin/bash
set -u -e -o pipefail -o verbose

go build -o blog
./blog
netlifyctl -A $NETLIFY_TOKEN deploy || true
cat netlifyctl-debug.log || true
