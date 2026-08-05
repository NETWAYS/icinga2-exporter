package collector

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/NETWAYS/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2InfluxDB2Collector struct {
	icingaClient icinga.IcingaClient
	logger       *slog.Logger
}

func NewIcinga2InfluxDB2Collector(client icinga.IcingaClient, logger *slog.Logger) *Icinga2InfluxDB2Collector {
	return &Icinga2InfluxDB2Collector{
		icingaClient: client,
		logger:       logger,
	}
}

func (collector *Icinga2InfluxDB2Collector) Describe(ch chan<- *prometheus.Desc) {
}

func (collector *Icinga2InfluxDB2Collector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointInfluxdb2Writer)

	if err != nil {
		collector.logger.Error("Could not retrieve InfluxDB2Writer metrics", "error", err.Error())
		return
	}

	for _, result := range perfdata.Results {
		for component, writers := range result.Status {
			for writer, metrics := range writers {
				for metricName, value := range metrics {
					// We only data about work/data items that are numeric values
					if !strings.HasPrefix(metricName, "data") && !strings.HasPrefix(metricName, "work") {
						continue
					}

					safeMetricName := strings.ReplaceAll(metricName, "-", "_")

					name := fmt.Sprintf("icinga2_%s_%s", component, safeMetricName)

					description := prometheus.NewDesc(name, "InfluxDB2Writer "+name, []string{"writer"}, nil)
					metric, err := prometheus.NewConstMetric(description, prometheus.GaugeValue, value, writer)

					if err != nil {
						collector.logger.Error("Error creating metric %s: %v", metricName, err)
						continue
					}

					ch <- metric
				}
			}
		}
	}
}
