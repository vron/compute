#/bin/bash
set -e

# First translate glsl shaders to spir-V
glslangValidator -G *.comp -o prog.spv
# glslangValidator -V *.comp -o prog.spv # This does not work either.

# Next convert this to llvm-IR
llvm-spirv -r prog.spv