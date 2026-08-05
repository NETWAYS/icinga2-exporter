# SPDX-License-Identifier: GPL-3.0-only

FROM docker.io/golang:alpine AS builder

ARG EXPORTER_VERSION=development
ARG EXPORTER_COMMIT=HEAD

WORKDIR /usr/local/src/exporter
COPY --chown=nobody:nogroup . .

RUN set -ex; \
    go build -ldflags="-s -w -X main.version=${EXPORTER_VERSION} -X main.commit=${EXPORTER_COMMIT}" -o /go/bin/icinga2-exporter

FROM docker.io/alpine:latest

RUN addgroup -S icinga_exporter && \
  adduser -S icinga_exporter -G icinga_exporter && \
  apk --no-cache add --update ca-certificates

COPY --from=builder /go/bin/icinga2-exporter /usr/sbin/icinga2-exporter

USER icinga_exporter
EXPOSE 9665

ENTRYPOINT ["/usr/sbin/icinga2-exporter"]
