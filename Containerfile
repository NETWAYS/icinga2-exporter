FROM docker.io/golang:alpine AS builder

WORKDIR /usr/local/src/exporter
COPY --chown=nobody:nogroup . .
RUN apk --no-cache add --update make && make build

FROM docker.io/alpine:latest

RUN addgroup -S icinga_exporter && \
  adduser -S icinga_exporter -G icinga_exporter && \
  apk --no-cache add --update ca-certificates

COPY --from=builder /usr/local/src/exporter/dist/icinga2-exporter /usr/sbin/icinga2-exporter

USER icinga_exporter
ENTRYPOINT ["/usr/sbin/icinga2-exporter"]

EXPOSE 9665
