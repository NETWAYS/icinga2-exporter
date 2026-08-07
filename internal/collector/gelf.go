package collector

import (
	"fmt"
	"log/slog"

	"github.com/NETWAYS/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2GELFCollector struct {
	icingaClient icinga.IcingaClient
	logger       *slog.Logger
}

func NewIcinga2GELFCollector(client icinga.IcingaClient, logger *slog.Logger) *Icinga2GELFCollector {
	return &Icinga2GELFCollector{
		icingaClient: client,
		logger:       logger,
	}
}

func (collector *Icinga2GELFCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (collector *Icinga2GELFCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointGELFWriter)

	if err != nil {
		collector.logger.Error("Could not retrieve GELFWriter metrics", "error", err.Error())
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

					description := prometheus.NewDesc(name, "GELFWriter "+name, []string{"writer"}, nil)
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
