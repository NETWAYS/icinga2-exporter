package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2CIBCollector struct {
	icingaClient          *icinga.Client
	logger                *slog.Logger
	uptime                *prometheus.Desc
	num_hosts_up          *prometheus.Desc
	num_hosts_down        *prometheus.Desc
	num_services_ok       *prometheus.Desc
	num_services_critical *prometheus.Desc
	avg_execution_time    *prometheus.Desc
	avg_latency           *prometheus.Desc
	max_execution_time    *prometheus.Desc
	max_latency           *prometheus.Desc
	min_execution_time    *prometheus.Desc
	min_latency           *prometheus.Desc
}

func NewIcinga2CIBCollector(client *icinga.Client, logger *slog.Logger) *Icinga2CIBCollector {
	return &Icinga2CIBCollector{
		icingaClient:          client,
		logger:                logger,
		uptime:                prometheus.NewDesc("icinga2_uptime", "Uptime of the instance", nil, nil),
		num_hosts_up:          prometheus.NewDesc("icinga2_num_hosts_up", "Number of Hosts Up", nil, nil),
		num_hosts_down:        prometheus.NewDesc("icinga2_num_hosts_down", "Number of Hosts Down", nil, nil),
		num_services_ok:       prometheus.NewDesc("icinga2_num_services_ok", "Number of Services OK", nil, nil),
		num_services_critical: prometheus.NewDesc("icinga2_num_services_critical", "Number of Services Critical", nil, nil),
		avg_execution_time:    prometheus.NewDesc("icinga2_avg_execution_time", "Average execution time", nil, nil),
		avg_latency:           prometheus.NewDesc("icinga2_avg_latency", "Average latency", nil, nil),
		max_execution_time:    prometheus.NewDesc("icinga2_max_execution_time", "Maximum execution time", nil, nil),
		max_latency:           prometheus.NewDesc("icinga2_max_latency", "Maximum latency", nil, nil),
		min_execution_time:    prometheus.NewDesc("icinga2_min_execution_time", "Minimum execution time", nil, nil),
		min_latency:           prometheus.NewDesc("icinga2_min_latency", "Minimum latency", nil, nil),
	}
}

func (collector *Icinga2CIBCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.uptime
	ch <- collector.num_hosts_up
	ch <- collector.num_hosts_down
	ch <- collector.num_services_ok
	ch <- collector.num_services_critical
	ch <- collector.avg_execution_time
	ch <- collector.avg_latency
	ch <- collector.max_execution_time
	ch <- collector.max_latency
	ch <- collector.min_execution_time
	ch <- collector.min_latency
}

func (collector *Icinga2CIBCollector) Collect(ch chan<- prometheus.Metric) {
	result, err := collector.icingaClient.GetCIBMetrics()

	if err != nil {
		collector.logger.Error("Could not retrieve CIB metrics", "error", err.Error())
		return
	}

	if len(result.Results) < 1 {
		collector.logger.Debug("No results for CIB metrics")
		return
	}

	r := result.Results[0]

	if v, ok := r.Status["uptime"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.uptime, prometheus.CounterValue, v)
	}
	if v, ok := r.Status["num_hosts_up"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_up, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_down"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_down, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_ok"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_ok, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_critical"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_critical, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["avg_execution_time"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.avg_execution_time, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["avg_latency"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.avg_latency, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["max_execution_time"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.max_execution_time, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["max_latency"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.max_latency, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["min_execution_time"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.min_execution_time, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["min_latency"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.min_latency, prometheus.GaugeValue, v)
	}
}
