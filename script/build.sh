#!/bin/bash
# TODO: move this over into a go file so it is cross platform...
set -e

SHADER=$1
CLANG=clang

# we build from scratch every time
mkdir -p build

# validate source
glslangValidator $SHADER
# copy the runtime to build folder
cp runtime/* build/
# transpile shader to cpp
(cd gl2c && env RUST_BACKTRACE=full cargo run -q -- ../$SHADER ../build/kernel.json)
# build the glue headers etc
go run glbind/*.go
# run the go file throug goimports
goimports -w build/kernel.go
# build the runtime to a static library
(cd build && clang++ lib.cpp -shared -o ./shader.so -fPIC -Ofast -Wall -Wextra -std=c++14 -fno-exceptions -ffast-math -Wno-unused-function -fno-math-errno) # TODO: enable the parameter again
#clang++ build/lib.cpp -O0 -shared -o build/shader.so -fPIE -std=c++14  # TODO: enable the parameter again
#ar rc build/shader.a build/lib.o
# move out only the relevant files to a subdir
mkdir -p build/go/build
cp build/*.go build/go
cp build/shader.so build/go/build
cp build/shader.so build/go
cp build/shared.h build/go/
# build the go pacae and run it to test it wors
(cd build/go && go test -v . )

