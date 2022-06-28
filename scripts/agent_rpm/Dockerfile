FROM docker.io/library/rockylinux@sha256:72afc2e1a20c9ddf56a81c51148ebcbe927c0a879849efe813bee77d69df1dd8

# Install rpm requirements
RUN dnf -y install git curl tar rpm-build rpmdevtools

# Install Golang
ARG GOLANG_VERSION=1.18
ARG GOLANG_ARCH=amd64
WORKDIR /usr/local
RUN curl -L https://dl.google.com/go/go${GOLANG_VERSION}.linux-${GOLANG_ARCH}.tar.gz | tar -xzf -
RUN ln -s /usr/local/go/bin/go /usr/bin/go
RUN ln -s /usr/local/go/bin/gofmt /usr/bin/gofmt

# Build RPM
WORKDIR /root/rpmbuild
ADD otto-agent.spec .
RUN mkdir SOURCES
COPY otto-*.tar.gz SOURCES/
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

VOLUME /root/rpmbuild/RPMS

ENTRYPOINT /entrypoint.sh
