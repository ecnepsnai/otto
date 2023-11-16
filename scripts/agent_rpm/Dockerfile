FROM docker.io/library/rockylinux:9

# Install rpm requirements
RUN dnf -y install git rpm-build rpmdevtools systemd-rpm-macros

# Install Golang
ARG GOLANG_VERSION=1.21.4
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
