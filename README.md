# Compute - static runtime for GLSL compute shaders

![Test](https://github.com/vron/compute/workflows/Test/badge.svg)

This project compiles GLSL 4.5/4.6 compute shaders to a C library
with low dependencies that can be called from any language allowing C bindings.
Additionally such a Go package is generated, providing a clean interface while using
cgo and the C library underneath.

Not all built-in functions in GLSL are yet implemented, this is a wip and bugs
are to be expected. Bug reports are appreciated, but please note
that the project main use case is for use in gioui and features not needed may not be
prioritized.

# Get started
For linux through docker:

1. Cd into the directory with your shader file(s) (*.comp)

2. Run:

    docker run -v $(pwd):/data vron/compute your_shader.comp

Unless you get errors (try with a simple shader first) you should have both the C and
the go library wrapper generated in the current directory.

For an example of using compute to generate a package, including how to call it
please see github.com/vron/computeexample

# Supported platforms
Currently supported on Linux, windows, osX amd64. Intended to support arm32 and arm64 as well.


# How to run locally
In order to run the project locally (not through docker) you will need a set of dependences,
in particular: recent clang, go, rust, lua, glslangValidator. For details please refer to the
Dockerfile to ensure you have all that is needed.

1. Cd into the top level directory of this project

2. Run (to run all tests and ensure everything works):

    go run test/main.go

3. Run (to build your shader):

    go run script/build.go path_to_your_shader.comp

4. Find your generated lib in ./build/go/*
