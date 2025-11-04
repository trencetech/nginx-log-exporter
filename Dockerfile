FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apk update && apk add --no-cache git bash

WORKDIR /nginx-log-exporter
COPY go.mod go.sum /nginx-log-exporter/
RUN go get ./...
ADD . /nginx-log-exporter

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -buildvcs=false -o nginx-log-exporter

FROM alpine:latest

COPY --from=builder /nginx-log-exporter/nginx-log-exporter /nginx-log-exporter

ARG BUILD_HASH
ENV BUILD_HASH=$BUILD_HASH

EXPOSE 9999
ENTRYPOINT ["/nginx-log-exporter"]
