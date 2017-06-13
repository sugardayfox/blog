#!/bin/bash
set -u -e -o pipefail

go run ./s/make-content-sha1.go
