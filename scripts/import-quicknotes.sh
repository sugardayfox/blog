#!/bin/bash
set -u -e -o pipefail

go run ./scripts/import-quicknotes.go
