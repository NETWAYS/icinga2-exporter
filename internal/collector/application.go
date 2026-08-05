package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2ApplicationCollector struct {
	icingaClient *icinga.Client
	logger       *slog.Logger
	info         *prometheus.Desc
}

func NewIcinga2ApplicationCollector(client *icinga.Client, logger *slog.Logger) *Icinga2ApplicationCollector {
	return &Icinga2ApplicationCollector{
		icingaClient: client,
		logger:       logger,
		info:         prometheus.NewDesc("icinga2_version_info", "A metric with a constant '1' value labeled by version", []string{"version"}, nil),
	}
}

func (collector *Icinga2ApplicationCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.info
}

func (collector *Icinga2ApplicationCollector) Collect(ch chan<- prometheus.Metric) {
	result, err := collector.icingaClient.GetApplicationMetrics()

	if err != nil {
		collector.logger.Error("Could not retrieve Application metrics", "error", err.Error())
		return
	}

	if len(result.Results) < 1 {
		collector.logger.Debug("No results for Application metrics")
		return
	}

	r := result.Results[0]

	ch <- prometheus.MustNewConstMetric(collector.info, prometheus.GaugeValue, 1, r.Status.IcingaApplication.App.Version)
}
