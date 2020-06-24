FROM ubuntu:focal
RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y wget gnupg2 software-properties-common git cmake build-essential glslang-tools clang-10

# Thet the llvm-spirv tool
RUN git clone -b llvm_release_100 https://github.com/KhronosGroup/SPIRV-LLVM-Translator.git && cd SPIRV-LLVM-Translator 
RUN mkdir SPIRV-LLVM-Translator/build && cd SPIRV-LLVM-Translator/build && cmake .. && make llvm-spirv -j`nproc`
RUN cp SPIRV-LLVM-Translator/build/tools/llvm-spirv/llvm-spirv /bin

ADD build.sh build.sh

VOLUME [ "/data" ]
WORKDIR /data

CMD /build.sh