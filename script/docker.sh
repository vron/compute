#!/bin/bash
# this script is needed so we can handle arguments sent to the container later

set -e
name=${1%.comp}
echo $1
echo "doing package $name"
ls -lah /data
go run script/build.go /data/$1
echo "done b"
cp -r /build/go/* /data/
sed -i "1,3s/kernel/$name/g" /data/kernel.go
sed -i "1,3s/kernel/$name/g" /data/types.go
mv /data/kernel.go /data/$name.go
