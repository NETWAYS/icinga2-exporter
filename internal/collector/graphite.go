package collector

import (
	"fmt"
	"log/slog"

	"github.com/NETWAYS/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2GraphiteCollector struct {
	icingaClient icinga.IcingaClient
	logger       *slog.Logger
}

func NewIcinga2GraphiteCollector(client icinga.IcingaClient, logger *slog.Logger) *Icinga2GraphiteCollector {
	return &Icinga2GraphiteCollector{
		icingaClient: client,
		logger:       logger,
	}
}

func (collector *Icinga2GraphiteCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (collector *Icinga2GraphiteCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointGraphiteWriter)

	if err != nil {
		collector.logger.Error("Could not retrieve GraphiteWriter metrics", "error", err.Error())
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

					description := prometheus.NewDesc(name, "GraphiteWriter "+name, []string{"writer"}, nil)
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
