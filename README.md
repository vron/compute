# Compute - static runtime for GLSL compute shaders

This project compiles (a subset of) vulan/openGL compute shaders to go packages
where the kernel will be executed on the CPU. The project uses 

This is very much a wor in progress and will only work for some compute shaders,
as driven by my own and GIO's (gioui.org) need. Use at your own risk, unsupported
shaders may fail in unexpected on unseen ways.

# Get started
Due to the complicated dependices the tool is easiest to use through docker. Put
your compute shader (e.g. kernel.comp) in a folder named "data" and from its
parent folder issue:

    docker run -v $(pwd)/data:/data vron/compute

If your shader is supported it should create the following files for you:

 - kernel.go
 - kernel.a
 - shared.h

That you should be able to use directly as a go package.

For an example of using compute to generate a package, including how to call it
please see github.com/vron/computeexample

# Alignment
Effectively std430 BUT vec3 and friends are rounded up to vec4 to allow for aligned sse2 operations.

# Testing
Testing is slow since it uses the full pipeling through the docker image for each
test case. To run all test cases run:

    make test

# Deficiencies
 - Currently only works for linux amd64 combination
 - Currently one thread only
 - Currently no synchronization primitives
 - Currently a very limited subset of GLSL is supported.
 - Do not support arrays of structs with arrays as children in input

# Architecture
TODO: Write once I have commited to a design


# TODO
 - add test for and ensure no name clash with our classes...
 - Create referenceGo implementations of all benchmars
 - Clean up go pacage
 - Rewor api such that array of array instead of merging all the arrays when that is the case..
 - support multi-file shader with macros
 - Gain more performance but letting the compiler now we have ensured all the alignements
 - Document that all fields not of slices will be COPIED to not mess with cgo ( i.e  vec4[1000] is a bad idea perforamnce wise..)