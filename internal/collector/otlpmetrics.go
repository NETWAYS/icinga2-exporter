package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2OTLPMetricsCollector struct {
	icingaClient                                        *icinga.Client
	logger                                              *slog.Logger
	otlpmetricswriter_otlp_metrics_work_queue_items     *prometheus.Desc
	otlpmetricswriter_otlp_metrics_work_queue_item_rate *prometheus.Desc
	otlpmetricswriter_otlp_metrics_data_buffer_items    *prometheus.Desc
	otlpmetricswriter_otlp_metrics_data_buffer_bytes    *prometheus.Desc
}

func NewIcinga2OTLPMetricsCollector(client *icinga.Client, logger *slog.Logger) *Icinga2OTLPMetricsCollector {
	return &Icinga2OTLPMetricsCollector{
		icingaClient: client,
		logger:       logger,
		otlpmetricswriter_otlp_metrics_work_queue_items:     prometheus.NewDesc("icinga2_otlpmetricswriter_otlp_metrics_work_queue_items", "OTLPMetricsWriter work queue items", nil, nil),
		otlpmetricswriter_otlp_metrics_work_queue_item_rate: prometheus.NewDesc("icinga2_otlpmetricswriter_otlp_metrics_work_queue_item_rate", "OTLPMetricsWriter work queue item rate", nil, nil),
		otlpmetricswriter_otlp_metrics_data_buffer_items:    prometheus.NewDesc("otlpmetricswriter_otlp_metrics_data_buffer_items", "OTLPMetricsWriter data buffer items", nil, nil),
		otlpmetricswriter_otlp_metrics_data_buffer_bytes:    prometheus.NewDesc("otlpmetricswriter_otlp_metrics_data_buffer_bytes", "OTLPMetricsWriter data buffer bytes", nil, nil),
	}
}

func (collector *Icinga2OTLPMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.otlpmetricswriter_otlp_metrics_work_queue_items
	ch <- collector.otlpmetricswriter_otlp_metrics_work_queue_item_rate
	ch <- collector.otlpmetricswriter_otlp_metrics_data_buffer_items
	ch <- collector.otlpmetricswriter_otlp_metrics_data_buffer_bytes
}

func (collector *Icinga2OTLPMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointGraphiteWriter)

	if err != nil {
		collector.logger.Error("Could not retrieve OTLPMetrics metrics", "error", err.Error())
		return
	}

	for _, datapoint := range perfdata {
		if datapoint.Label == "otlpmetricswriter_otlp_metrics_work_queue_items" {
			ch <- prometheus.MustNewConstMetric(collector.otlpmetricswriter_otlp_metrics_work_queue_items, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "otlpmetricswriter_otlp_metrics_work_queue_item_rate" {
			ch <- prometheus.MustNewConstMetric(collector.otlpmetricswriter_otlp_metrics_work_queue_item_rate, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "otlpmetricswriter_otlp_metrics_data_buffer_items" {
			ch <- prometheus.MustNewConstMetric(collector.otlpmetricswriter_otlp_metrics_data_buffer_items, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "otlpmetricswriter_otlp_metrics_data_buffer_bytes" {
			ch <- prometheus.MustNewConstMetric(collector.otlpmetricswriter_otlp_metrics_data_buffer_bytes, prometheus.GaugeValue, datapoint.Value)
		}
	}
}
