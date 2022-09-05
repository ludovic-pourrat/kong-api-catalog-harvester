FROM golang:1.19.0 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /api-catalog-harvester

FROM kong

USER root
COPY --from=build /api-catalog-harvester /usr/local/bin/
USER kong
