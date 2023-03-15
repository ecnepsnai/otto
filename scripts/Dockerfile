ARG ALPINE_HASH
FROM docker.io/library/alpine@sha256:${ALPINE_HASH}
LABEL maintainer="Ian Spence <ian@ecnepsnai.com>"
LABEL org.opencontainers.image.authors="Ian Spence <ian@ecnepsnai.com>"
LABEL org.opencontainers.image.source=https://github.com/ecnepsnai/otto
LABEL org.opencontainers.image.title="otto"
LABEL org.opencontainers.image.description="Otto Server"

EXPOSE 8080
RUN mkdir /otto_data
VOLUME [ "/otto_data" ]
COPY "otto" /otto

ENTRYPOINT [ "/otto/otto", "--data-dir", "/otto_data", "-b", "0.0.0.0:8080" ]
