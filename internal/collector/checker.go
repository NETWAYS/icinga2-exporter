package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2CheckerCollector struct {
	icingaClient                     *icinga.Client
	logger                           *slog.Logger
	checkercomponent_checker_idle    *prometheus.Desc
	checkercomponent_checker_pending *prometheus.Desc
}

func NewIcinga2CheckerCollector(client *icinga.Client, logger *slog.Logger) *Icinga2CheckerCollector {
	return &Icinga2CheckerCollector{
		icingaClient:                     client,
		logger:                           logger,
		checkercomponent_checker_idle:    prometheus.NewDesc("icinga2_checkercomponent_checker_idle", "CheckerComponent idle", nil, nil),
		checkercomponent_checker_pending: prometheus.NewDesc("icinga2_checkercomponent_checker_pending", "CheckerComponent pending", nil, nil),
	}
}

func (collector *Icinga2CheckerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.checkercomponent_checker_idle
	ch <- collector.checkercomponent_checker_pending
}

func (collector *Icinga2CheckerCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointCheckerComponent)

	if err != nil {
		collector.logger.Error("Could not retrieve CheckerComponent metrics", "error", err.Error())
		return
	}

	for _, datapoint := range perfdata {
		if datapoint.Label == "checkercomponent_checker_idle" {
			ch <- prometheus.MustNewConstMetric(collector.checkercomponent_checker_idle, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "checkercomponent_checker_pending" {
			ch <- prometheus.MustNewConstMetric(collector.checkercomponent_checker_pending, prometheus.GaugeValue, datapoint.Value)
		}
	}
}
