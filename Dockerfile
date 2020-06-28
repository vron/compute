FROM ubuntu:focal
# TODO: update to clang 11
RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y build-essential glslang-tools clang golang-go cargo rustc

#RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
#RUN echo 'source $HOME/.cargo/env' >> $HOME/.bashrc

RUN mkdir /data

COPY gl2c gl2c
RUN cd gl2c && env PATH="$HOME/.cargo/bin:$PATH" cargo build
COPY glbind glbind
COPY runtime runtime
COPY script script

VOLUME [ "/build" ]

CMD /script/docker.sh