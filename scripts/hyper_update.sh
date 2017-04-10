#!/bin/bash
set -u -e -o pipefail

hyper fip detach blog
hyper stop blog
hyper rm blog
hyper run --size=s2 --restart=unless-stopped -d -p 80 -v blog:/data --name blog blog:latest
hyper fip attach 209.177.93.141 blog
