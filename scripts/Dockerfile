FROM alpine

RUN echo http://dl-2.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories
RUN apk --no-cache add --update shadow

RUN addgroup -S otto && adduser -S otto -G otto

EXPOSE 8080
RUN mkdir /otto_data
VOLUME [ "/otto_data" ]
ADD "otto" /otto
ADD "entrypoint.sh" /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]