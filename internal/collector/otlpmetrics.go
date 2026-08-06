package collector

import (
	"fmt"
	"log/slog"

	"github.com/NETWAYS/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2OTLPMetricsCollector struct {
	icingaClient icinga.IcingaClient
	logger       *slog.Logger
}

func NewIcinga2OTLPMetricsCollector(client icinga.IcingaClient, logger *slog.Logger) *Icinga2OTLPMetricsCollector {
	return &Icinga2OTLPMetricsCollector{
		icingaClient: client,
		logger:       logger,
	}
}

func (collector *Icinga2OTLPMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (collector *Icinga2OTLPMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointOTLPMetricsWriter)

	if err != nil {
		collector.logger.Error("Could not retrieve OTLPMetricsWriter metrics", "error", err.Error())
		return
	}

	for _, result := range perfdata.Results {
		for component, writers := range result.Status {
			for writer, metrics := range writers {
				for metricName, value := range metrics {
					if !isPerfdataMetric(metricName) {
						continue
					}

					safeMetricName := ensureValidMetricName(metricName)

					name := fmt.Sprintf("icinga2_%s_%s", component, safeMetricName)

					description := prometheus.NewDesc(name, "OTLPMetricsWriter "+name, []string{"writer"}, nil)
					metric, err := prometheus.NewConstMetric(description, prometheus.GaugeValue, value, writer)

					if err != nil {
						collector.logger.Error("Error creating metric "+metricName, "error", err.Error())
						continue
					}

					ch <- metric
				}
			}
		}
	}
}
