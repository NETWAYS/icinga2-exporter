package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2InfluxDBCollector struct {
	icingaClient                                 *icinga.Client
	logger                                       *slog.Logger
	influxdbwriter_influxdb_work_queue_items     *prometheus.Desc
	influxdbwriter_influxdb_work_queue_item_rate *prometheus.Desc
	influxdbwriter_influxdb_data_queue_items     *prometheus.Desc
}

func NewIcinga2InfluxDBCollector(client *icinga.Client, logger *slog.Logger) *Icinga2InfluxDBCollector {
	return &Icinga2InfluxDBCollector{
		icingaClient:                             client,
		logger:                                   logger,
		influxdbwriter_influxdb_work_queue_items: prometheus.NewDesc("icinga2_influxdbwriter_influxdb_work_queue_items", "InfluxDBWriter work queue items", nil, nil),
		influxdbwriter_influxdb_work_queue_item_rate: prometheus.NewDesc("icinga2_influxdbwriter_influxdb_work_queue_item_rate", "InfluxDBWriter work queue item rate", nil, nil),
		influxdbwriter_influxdb_data_queue_items:     prometheus.NewDesc("icinga2_influxdbwriter_influxdb_data_queue_items", "InfluxDBWriter data queue items", nil, nil),
	}
}

func (collector *Icinga2InfluxDBCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.influxdbwriter_influxdb_work_queue_items
	ch <- collector.influxdbwriter_influxdb_work_queue_item_rate
	ch <- collector.influxdbwriter_influxdb_data_queue_items
}

func (collector *Icinga2InfluxDBCollector) Collect(ch chan<- prometheus.Metric) {
	perfdata, err := collector.icingaClient.GetPerfdataMetrics(icinga.EndpointInfluxdbWriter)

	if err != nil {
		collector.logger.Error("Could not retrieve InfluxDB metrics", "error", err.Error())
		return
	}

	for _, datapoint := range perfdata {
		if datapoint.Label == "influxdbwriter_influxdb_work_queue_items" {
			ch <- prometheus.MustNewConstMetric(collector.influxdbwriter_influxdb_work_queue_items, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "influxdbwriter_influxdb_work_queue_item_rate" {
			ch <- prometheus.MustNewConstMetric(collector.influxdbwriter_influxdb_work_queue_item_rate, prometheus.GaugeValue, datapoint.Value)
		}

		if datapoint.Label == "influxdbwriter_influxdb_data_queue_items" {
			ch <- prometheus.MustNewConstMetric(collector.influxdbwriter_influxdb_data_queue_items, prometheus.GaugeValue, datapoint.Value)
		}
	}
}
