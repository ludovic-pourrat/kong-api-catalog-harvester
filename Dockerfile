FROM golang:1.19.0 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY utils/*.go ./utils/
COPY utils/pathtrie/*.go ./utils/pathtrie/
COPY types/*.go ./types/
COPY factories/*.go ./factories/
COPY *.go ./

RUN CGO_ENABLED=0 go build -o /api-catalog-harvester

FROM kong:2.8

USER root
COPY --from=build /api-catalog-harvester /usr/local/bin/
USER kong

