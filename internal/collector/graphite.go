package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2GraphiteCollector struct {
	icingaClient                                 *icinga.Client
	logger                                       *slog.Logger
	graphitewriter_graphite_work_queue_items     *prometheus.Desc
	graphitewriter_graphite_work_queue_item_rate *prometheus.Desc
	graphitewriter_graphite_data_queue_items     *prometheus.Desc
}

func NewIcinga2GraphiteCollector(client *icinga.Client, logger *slog.Logger) *Icinga2GraphiteCollector {
	return &Icinga2GraphiteCollector{
		icingaClient:                             client,
		logger:                                   logger,
		graphitewriter_graphite_work_queue_items: prometheus.NewDesc("icinga2_graphitewriter_graphite_work_queue_items", "GraphiteWriter work queue items", nil, nil),
		graphitewriter_graphite_work_queue_item_rate: prometheus.NewDesc("icinga2_graphitewriter_graphite_work_queue_item_rate", "GraphiteWriter work queue item rate", nil, nil),
	}
}

func (collector *Icinga2GraphiteCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.graphitewriter_graphite_work_queue_items
	ch <- collector.graphitewriter_graphite_work_queue_item_rate
}

func (collector *Icinga2GraphiteCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointGraphiteWriter)

	if err != nil {
		collector.logger.Error("Could not retrieve Graphite metrics", "error", err.Error())
		return
	}

	for _, datapoint := range perfdata {
		if datapoint.Label == "graphitewriter_graphite_work_queue_items" {
			ch <- prometheus.MustNewConstMetric(collector.graphitewriter_graphite_work_queue_items, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "graphitewriter_graphite_work_queue_item_rate" {
			ch <- prometheus.MustNewConstMetric(collector.graphitewriter_graphite_work_queue_item_rate, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "graphitewriter_graphite_data_queue_items" {
			ch <- prometheus.MustNewConstMetric(collector.graphitewriter_graphite_data_queue_items, prometheus.GaugeValue, datapoint.Value)
		}
	}
}
