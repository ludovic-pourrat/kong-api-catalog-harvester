FROM node:alpine3.16@sha256:f08168de449131d96a16a9c042f96dc3169678907f120eee8d5ecc10ca75bb48

COPY entrypoint.sh /entrypoint.sh

RUN \
  addgroup tool && \
  adduser -G tool -s /bin/sh -D tool && \
  apk update && apk add coreutils curl && rm -rf /var/cache/apk/* && \
  chmod u=rwx,g=rx,o= /entrypoint.sh

RUN npm install -g newman

WORKDIR /etc/newman

RUN chown -R tool:root /etc/newman /entrypoint.sh

USER tool

ENTRYPOINT ["/entrypoint.sh"]