#!/bin/bash
# this script is needed so we can handle arguments sent to the container later

set -e
name=${1%.comp}
echo "doing package $name"
go run script/build.go /data/$1
cp -r /build/go/* /data/
sed -i "1,3s/kernel/$name/g" /data/kernel.go
mv /data/kernel.go /data/$name.go
