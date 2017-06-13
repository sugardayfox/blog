#!/bin/bash
set -u -e -o pipefail

go tool vet -printfuncs=httpErrorf:1,panicif:1,Noticef,Errorf .
