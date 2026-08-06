package collector

import (
	"fmt"
	"log/slog"

	"github.com/NETWAYS/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2InfluxDBCollector struct {
	icingaClient icinga.IcingaClient
	logger       *slog.Logger
}

func NewIcinga2InfluxDBCollector(client icinga.IcingaClient, logger *slog.Logger) *Icinga2InfluxDBCollector {
	return &Icinga2InfluxDBCollector{
		icingaClient: client,
		logger:       logger,
	}
}

func (collector *Icinga2InfluxDBCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (collector *Icinga2InfluxDBCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointInfluxdbWriter)

	if err != nil {
		collector.logger.Error("Could not retrieve InfluxDBWriter metrics", "error", err.Error())
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

					description := prometheus.NewDesc(name, "InfluxDBWriter "+name, []string{"writer"}, nil)
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
