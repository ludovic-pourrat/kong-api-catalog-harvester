#
# Prism - See https://stoplight.io/open-source/prism
#
#
#
FROM node:alpine3.16@sha256:f08168de449131d96a16a9c042f96dc3169678907f120eee8d5ecc10ca75bb48

COPY entrypoint.sh /entrypoint.sh

RUN \
  addgroup mock && \
  adduser -G mock -s /bin/sh -D mock && \
  apk update && apk add coreutils && rm -rf /var/cache/apk/* && \
  chmod u=rwx,g=rx,o= /entrypoint.sh

RUN npm install -g @stoplight/prism-cli

WORKDIR /workspace

RUN chown -R mock:root /workspace /entrypoint.sh

USER mock

EXPOSE 4010

ENTRYPOINT ["/entrypoint.sh"]
