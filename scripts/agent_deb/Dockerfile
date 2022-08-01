FROM docker.io/library/debian:stable
RUN apt-get -y update && apt-get -y install curl

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

VOLUME [ "/ottoagent" ]
WORKDIR /

ENTRYPOINT [ "/entrypoint.sh" ]