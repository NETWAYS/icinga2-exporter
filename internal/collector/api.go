package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2APICollector struct {
	icingaClient               *icinga.Client
	logger                     *slog.Logger
	api_num_conn_endpoints     *prometheus.Desc
	api_num_not_conn_endpoints *prometheus.Desc
	api_num_endpoints          *prometheus.Desc
	api_num_http_clients       *prometheus.Desc
}

func NewIcinga2APICollector(client *icinga.Client, logger *slog.Logger) *Icinga2APICollector {
	return &Icinga2APICollector{
		icingaClient:               client,
		logger:                     logger,
		api_num_conn_endpoints:     prometheus.NewDesc("icinga2_api_num_conn_endpoints", "Number of connected Endpoints", nil, nil),
		api_num_endpoints:          prometheus.NewDesc("icinga2_api_num_endpoints", "Number of Endpoints", nil, nil),
		api_num_not_conn_endpoints: prometheus.NewDesc("icinga2_api_num_not_conn_endpoints", "Number of not connected Endpoints", nil, nil),
		api_num_http_clients:       prometheus.NewDesc("icinga2_api_num_http_clients", "Number of HTTP Clients", nil, nil),
	}
}

func (collector *Icinga2APICollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.api_num_conn_endpoints
	ch <- collector.api_num_not_conn_endpoints
	ch <- collector.api_num_endpoints
	ch <- collector.api_num_http_clients
}

func (collector *Icinga2APICollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointApiListener)

	if err != nil {
		collector.logger.Error("Could not retrieve ApiListener metrics", "error", err.Error())
		return
	}

	for _, datapoint := range perfdata {
		if datapoint.Label == "api_num_conn_endpoints" {
			ch <- prometheus.MustNewConstMetric(collector.api_num_conn_endpoints, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "api_num_not_conn_endpoints" {
			ch <- prometheus.MustNewConstMetric(collector.api_num_conn_endpoints, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "api_num_endpoints" {
			ch <- prometheus.MustNewConstMetric(collector.api_num_conn_endpoints, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "api_num_http_clients" {
			ch <- prometheus.MustNewConstMetric(collector.api_num_conn_endpoints, prometheus.GaugeValue, datapoint.Value)
		}
	}
}
