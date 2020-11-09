FROM alpine

EXPOSE 8080
RUN mkdir /otto_data
VOLUME [ "/otto_data" ]
COPY "otto" /otto
ADD "entrypoint.sh" /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
