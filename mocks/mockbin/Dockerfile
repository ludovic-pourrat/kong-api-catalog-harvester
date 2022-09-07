#
# Mockbin - See https://github.com/Kong/mockbin
#
#
#
FROM node:14.20.0-alpine3.15@sha256:f00e9497a23928d5cef5d6f09001e5c8072be9409d63c68231dfb74489d2124e

COPY entrypoint.sh /entrypoint.sh

RUN \
  addgroup mock && \
  adduser -G mock -s /bin/sh -D mock && \
  apk update && apk add coreutils && rm -rf /var/cache/apk/* && \
  chmod u=rwx,g=rx,o= /entrypoint.sh

RUN npm install -g mockbin

WORKDIR /workspace

RUN chown -R mock:root /workspace /entrypoint.sh

USER mock

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
