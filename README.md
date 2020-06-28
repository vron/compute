# Compute - static runtime for GLSL compute shaders

This project compiles (a subset of) vulan/openGL compute shaders to go packages
where the kernel will be executed on the CPU. The project uses 

This is very much a wor in progress and will only work for some compute shaders,
as driven by my own and GIO's (gioui.org) need. Use at your own risk, unsupported
shaders may fail in unexpected on unseen ways.

# Get started
Due to the complicated dependices the tool is easiest to use through docker. Put
your compute shader (e.g. kernel.comp) in a folder and inside that folder issue:

    docker run vron/compute kernel.comp

If your shader is supported it should create the following files for you:

 - kernel.go
 - kernel.a
 - shared.h

That you should be able to use directly as a go package.

For an example of using compute to generate a package, including how to call it
please see github.com/vron/computeexample

# Testing
Testing is slow since it uses the full pipeling through the docker image for each
test case. To run all test cases run:

    make test

# Deficiencies
 - Currently only works for linux amd64 combination
 - Currently one thread only
 - Currently no synchronization primitives
 - Currently a very limited subset of GLSL is supported.

# Architecture
TODO: Write once I have commited to a design


# TODO
 - Ensure e.g. int sizes match when on 32 bit platform - likely not (opencl - platform that is)
 - Ensure struct alignments similarity between go and c - likely a problem? (or does cgo handle)
 - support multi-file shader