#!/bin/bash
set -u -e -o pipefail

go run ./s/netlify_build.go
cd ./netlify_static
echo "About to deploy"
#netlifyctl deploy
