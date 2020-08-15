FROM ubuntu:focal
RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y build-essential glslang-tools clang golang-go cargo rustc git
RUN go get golang.org/x/tools/cmd/goimports && cp /root/go/bin/goimports /bin/ && \
    mkdir /data

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y lua5.3

# Prebuild rust tool to make it go quicker for user
COPY gl2c gl2c
RUN cd gl2c && env PATH="$HOME/.cargo/bin:$PATH" cargo build



COPY glbind glbind
COPY go.mod go.mod
COPY go.sum go.sum
# Pre-download go modules to make it quicker for user
RUN go get github.com/termie/go-shutil 

COPY runtime runtime
COPY script script

VOLUME [ "/build" ]

ENTRYPOINT ["/script/docker.sh"]