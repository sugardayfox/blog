#!/bin/bash
set -u -e -o pipefail

go run ./scripts/make-content-sha1.go
