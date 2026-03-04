package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2InfluxDB2Collector struct {
	icingaClient                                   *icinga.Client
	logger                                         *slog.Logger
	influxdb2writer_influxdb2_work_queue_items     *prometheus.Desc
	influxdb2writer_influxdb2_work_queue_item_rate *prometheus.Desc
	influxdb2writer_influxdb2_data_queue_items     *prometheus.Desc
}

func NewIcinga2InfluxDB2Collector(client *icinga.Client, logger *slog.Logger) *Icinga2InfluxDB2Collector {
	return &Icinga2InfluxDB2Collector{
		icingaClient: client,
		logger:       logger,
		influxdb2writer_influxdb2_work_queue_items:     prometheus.NewDesc("icinga2_influxdb2writer_influxdb2_work_queue_items", "InfluxDB2Writer work queue items", nil, nil),
		influxdb2writer_influxdb2_work_queue_item_rate: prometheus.NewDesc("icinga2_influxdb2writer_influxdb2_work_queue_item_rate", "InfluxDB2Writer work queue item rate", nil, nil),
		influxdb2writer_influxdb2_data_queue_items:     prometheus.NewDesc("icinga2_influxdb2writer_influxdb2_data_queue_items", "InfluxDB2Writer data queue items", nil, nil),
	}
}

func (collector *Icinga2InfluxDB2Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.influxdb2writer_influxdb2_work_queue_items
	ch <- collector.influxdb2writer_influxdb2_work_queue_item_rate
	ch <- collector.influxdb2writer_influxdb2_data_queue_items
}

func (collector *Icinga2InfluxDB2Collector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointInfluxdb2Writer)

	if err != nil {
		collector.logger.Error("Could not retrieve InfluxDB2 metrics", "error", err.Error())
		return
	}

	for _, datapoint := range perfdata {
		if datapoint.Label == "influxdb2writer_influxdb2_work_queue_items" {
			ch <- prometheus.MustNewConstMetric(collector.influxdb2writer_influxdb2_work_queue_items, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "influxdb2writer_influxdb2_work_queue_item_rate" {
			ch <- prometheus.MustNewConstMetric(collector.influxdb2writer_influxdb2_work_queue_item_rate, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "influxdb2writer_influxdb2_data_queue_items" {
			ch <- prometheus.MustNewConstMetric(collector.influxdb2writer_influxdb2_data_queue_items, prometheus.GaugeValue, datapoint.Value)
		}
	}
}
