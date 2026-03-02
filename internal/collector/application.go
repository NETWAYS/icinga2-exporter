package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2ApplicationCollector struct {
	icingaClient *icinga.Client
	logger       *slog.Logger
	info         *prometheus.GaugeVec
}

func NewIcinga2ApplicationCollector(client *icinga.Client, logger *slog.Logger) *Icinga2ApplicationCollector {
	return &Icinga2ApplicationCollector{
		icingaClient: client,
		logger:       logger,
		info: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "icinga2",
				Name:      "version_info",
				Help:      "A metric with a constant '1' value labeled by version",
			},
			[]string{"version"},
		),
	}
}

func (collector *Icinga2ApplicationCollector) Describe(ch chan<- *prometheus.Desc) {
	collector.info.Describe(ch)
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

	// TODO: Use a custom unmarshal to avoid this
	r := result.Results[0]

	collector.info.Reset()

	collector.info.With(prometheus.Labels{
		"version": r.Status.IcingaApplication.App.Version,
	}).Set(1)

	collector.info.Collect(ch)
}
