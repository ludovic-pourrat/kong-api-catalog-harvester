FROM golang:1.19.0 as build

WORKDIR /app

RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY utils/*.go ./utils/
COPY types/*.go ./types/
COPY factories/*.go ./factories/
COPY *.go ./

RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /api-catalog-harvester

FROM kong:2.8

USER root
COPY debug/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY --from=build /go/bin/dlv /
COPY --from=build /api-catalog-harvester /usr/local/bin/
USER kong
