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
	result, err := collector.icingaClient.GetCheckerComponentMetrics()

	if err != nil {
		collector.logger.Error("Could not retrieve CheckerComponent metrics", "error", err.Error())
		return
	}

	if len(result.Results) < 1 {
		collector.logger.Debug("No results for CheckerComponent metrics")
		return
	}

	r := result.Results[0]
	// There might be a better way
	var perfdata = make(map[string]float64, len(r.Perfdata))
	for _, v := range r.Perfdata {
		perfdata[v.Label] = v.Value
	}

	if v, ok := perfdata["checkercomponent_checker_idle"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.checkercomponent_checker_idle, prometheus.GaugeValue, v)
	}

	if v, ok := perfdata["checkercomponent_checker_pending"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.checkercomponent_checker_pending, prometheus.GaugeValue, v)
	}
}
