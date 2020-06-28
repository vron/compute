FROM ubuntu:focal
# TODO: update to clang 11
RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y wget gnupg2 software-properties-common git cmake build-essential glslang-tools clang golang-go

# Thet the llvm-spirv tool
#RUN git clone -b llvm_release_100 https://github.com/KhronosGroup/SPIRV-LLVM-Translator.git && cd SPIRV-LLVM-Translator 
#RUN mkdir SPIRV-LLVM-Translator/build && cd SPIRV-LLVM-Translator/build && cmake .. && make llvm-spirv -j`nproc`
#RUN cp SPIRV-LLVM-Translator/build/tools/llvm-spirv/llvm-spirv /bin

ADD makefile /makefile
COPY src src
ADD src/docker.sh /docker.sh
RUN mkdir /data

VOLUME [ "/data" ]

CMD /docker.sh