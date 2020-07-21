#!/bin/bash
# this script is needed so we can handle arguments sent to the container later

set -e
go run script/build.go /data/$1
cp -r /build/go/* /data/