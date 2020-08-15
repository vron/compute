#!/bin/bash
set -e
go run script/clean.go
go run script/build.go test/leaks/shader.comp
cp test/leaks/main.cpp build/main.cpp


(cd build && clang++ -std=c++2a -fvisibility=hidden -Wall -Wextra -Werror -Wno-unused-function -Ofast -ffast-math -fno-math-errno  lib.cpp main.cpp co/arch/*.S)
