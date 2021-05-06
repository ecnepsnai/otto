FROM arm64v8/alpine:3
LABEL maintainer="Ian Spence <ian@ecn.io>"
LABEL org.opencontainers.image.source=https://github.com/ecnepsnai/otto

EXPOSE 8080
RUN mkdir /otto_data
VOLUME [ "/otto_data" ]
COPY "otto" /otto
ADD "entrypoint.sh" /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
