# Icinga2 API exporter

Prometheus exporter for the Icinga2 API.

## Installation and Usage

The `icinga2_exporter` listens on HTTP port 9665 by default.
See the `-help` output for more options.

```
-collector.apilistener
      Include APIListener data
-collector.cib
      Include CIB data
-debug
      Enable debug logging
-icinga.api string
      Path to the Icinga2 API (default "https://localhost:5665/v1")
-icinga.cafile string
      Path to the Icinga2 API TLS CA
-icinga.certfile string
      Path to the Icinga2 API TLS cert
-icinga.insecure
      Skip TLS verification for Icinga2 API
-icinga.keyfile string
      Path to the Icinga2 API TLS key
-icinga.password string
      Icinga2 API Password
-icinga.username string
      Icinga2 API Username
-version
      Print version
-web.cache-ttl uint
      Cache lifetime in seconds for the Icinga API responses (default 60)
-web.listen-address string
      Address on which to expose metrics and web interface. (default ":9665")
-web.metrics-path string
      Path under which to expose metrics. (default "/metrics")
```

## Collectors

By default only the `IcingaApplication` metrics of the status API are collected.

There are more collectors that can be activated via the CLI.
The tables below list all existing collectors.

| Collector     | Flag       |
| ------------- | ---------- |
| APIListener   | `-collector.apilistener` |
| CIB           | `-collector.cib`         |

# Development

Running tests:

```
make test
make coverage
```
