# Compute - static runtime for GLSL compute shaders

![Test](https://github.com/vron/compute/workflows/Test/badge.svg)

This project compiles (a subset of) Vulkan/OpenGL compute shaders to a C library
with low dependencies that can be called from any language allowing C bindings.
Additionally such a Go package is generated, providing a clean interface while using
cgo and the C library underneath.

Not all built-in functions in GLSL are yet implemented, and although there is some
test coverage far from everything (in particular related to synchronization primitives)
is as of yet tested. Bug reports with minimal examples are appreciated, but please note
that the project main use case is for use in gioui and features not needed may not be
prioritized.

# Get started
Due to the complicated dependencies the tool is easiest to use through docker. If you
want to run and use it locally please refer to the source code. Note that Windows and macOS
are supported, but not through the docker file.

To simply use this project to generate a go package from a compute shader such that you
can import and use it directly, the easiest is to run it through docker (only supported
on Linux for now, see TODO below, also the image is huge - could probably be slimmed down
a bit, but ut must include llvm, rust, go etc. etc.):

1. Cd into the directory with your shader file (*.comp file)

2. Run:

    docker run -v $(pwd):/data vron/compute your_shader.comp

Unless you get errors (try with a simple shader first) you should have both the C and
the go library wrapper generated.

For an example of using compute to generate a package, including how to call it
please see github.com/vron/computeexample

# Known limitations
 - Currently only works on x64
 - Several built in glsl functions missing

# Development get-started
There has not been time to write a detailed explanation, but the best starting point to understand
what is going on is to start from the top level command:

    go run tests/main.go

which runs all the tests and follow what it does (see also scripts/build.go).

# TODO
 - Cross compile for all platforms from the Docker image.
 - Implement support for ARM64, Android, Linux, Windows
 - Create referenceGo implementations of all benchmarks
 - Add test cases for name collisions and fix the naming...
 - Support multi-file shader with macros
 - Gain more performance but letting the compiler now we have ensured all the alignments
 - Add pathological test cases such as wg size 0 etc. etc.
 - Set the NDEBUG flag to remove all those asserts..    
 - Document the alignment for []byte in both the C header and the go package.
 - Check if we have any problems with false sharing
 - When running the tests, first pre-build the headers to save time on each?