#!/bin/bash
set -u -e -o pipefail

pyflakes scripts/*py
