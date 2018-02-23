#!/bin/bash
set -u -e -o pipefail -o verbose

go get -v -u github.com/netlify/netlifyctl
go get -v -u github.com/golang/dep/cmd/dep
dep ensure -v
