#!/bin/sh

sysctl -w fs.file-max=128000

l=`cat /proc/sys/fs/file-max`
echo "setting open files limit (ulimit -n) to {$l}"
ulimit -Hn ${l}
ulimit -Sn ${l}
echo "ulimit -Sn:"
ulimit -Sn

./blog -addr=:80
