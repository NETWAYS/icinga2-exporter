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
	result, err := collector.icingaClient.GetApiListenerMetrics()

	if err != nil {
		collector.logger.Error("Could not retrieve ApiListener metrics", "error", err.Error())
		return
	}

	if len(result.Results) < 1 {
		collector.logger.Debug("No results for ApiListener metrics")
		return
	}

	r := result.Results[0]
	// There might be a better way
	var perfdata = make(map[string]float64, len(r.Perfdata))
	for _, v := range r.Perfdata {
		perfdata[v.Label] = v.Value
	}

	ch <- prometheus.MustNewConstMetric(collector.api_num_conn_endpoints, prometheus.GaugeValue, perfdata["api_num_conn_endpoints"])
	ch <- prometheus.MustNewConstMetric(collector.api_num_not_conn_endpoints, prometheus.GaugeValue, perfdata["api_num_not_conn_endpoints"])
	ch <- prometheus.MustNewConstMetric(collector.api_num_endpoints, prometheus.GaugeValue, perfdata["api_num_endpoints"])
	ch <- prometheus.MustNewConstMetric(collector.api_num_http_clients, prometheus.GaugeValue, perfdata["api_num_http_clients"])
}
