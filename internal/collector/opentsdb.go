package collector

import (
	"fmt"
	"log/slog"

	"github.com/NETWAYS/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2OpenTSDBCollector struct {
	icingaClient icinga.IcingaClient
	logger       *slog.Logger
}

func NewIcinga2OpenTSDBCollector(client icinga.IcingaClient, logger *slog.Logger) *Icinga2OpenTSDBCollector {
	return &Icinga2OpenTSDBCollector{
		icingaClient: client,
		logger:       logger,
	}
}

func (collector *Icinga2OpenTSDBCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (collector *Icinga2OpenTSDBCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointOpenTsdbWriter)

	if err != nil {
		collector.logger.Error("Could not retrieve OpenTSDBWriter metrics", "error", err.Error())
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

					description := prometheus.NewDesc(name, "OpenTSDBWriter "+name, []string{"writer"}, nil)
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
