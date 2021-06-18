FROM alpine:latest

RUN mkdir /otto
ADD entrypoint.sh /entrypoint.sh
ADD otto /otto/otto
RUN chmod a+x /entrypoint.sh
RUN chmod a+x /otto/otto

EXPOSE 12444

ENTRYPOINT [ "/entrypoint.sh" ]