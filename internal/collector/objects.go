package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	zonesObject = "zones"
	usersObject = "users"
)

type Icinga2ObjectsCollector struct {
	icingaClient      *icinga.Client
	logger            *slog.Logger
	objects_zones_num *prometheus.Desc
	objects_users_num *prometheus.Desc
}

func NewIcinga2ObjectsCollector(client *icinga.Client, logger *slog.Logger) *Icinga2ObjectsCollector {
	return &Icinga2ObjectsCollector{
		icingaClient:      client,
		logger:            logger,
		objects_zones_num: prometheus.NewDesc("icinga2_objects_zones_num", "Number of Zone objects", nil, nil),
		objects_users_num: prometheus.NewDesc("icinga2_objects_users_num", "Number of User objects", nil, nil),
	}
}

func (collector *Icinga2ObjectsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.objects_zones_num
	ch <- collector.objects_users_num
}

func (collector *Icinga2ObjectsCollector) Collect(ch chan<- prometheus.Metric) {
	collector.collectZones(ch)
	collector.collectUsers(ch)
}

func (collector *Icinga2ObjectsCollector) collectZones(ch chan<- prometheus.Metric) {
	result, err := collector.icingaClient.GetObjects(zonesObject)

	if err != nil {
		collector.logger.Error("Could not retrieve Zones metrics", "error", err.Error())
		return
	}

	if len(result.Results) < 1 {
		collector.logger.Debug("No results for Zones metrics")
		return
	}

	zoneCount := float64(len(result.Results))

	ch <- prometheus.MustNewConstMetric(collector.objects_zones_num, prometheus.GaugeValue, zoneCount)
}

func (collector *Icinga2ObjectsCollector) collectUsers(ch chan<- prometheus.Metric) {
	result, err := collector.icingaClient.GetObjects(usersObject)

	if err != nil {
		collector.logger.Error("Could not retrieve Users metrics", "error", err.Error())
		return
	}

	if len(result.Results) < 1 {
		collector.logger.Debug("No results for Users metrics")
		return
	}

	userCount := float64(len(result.Results))

	ch <- prometheus.MustNewConstMetric(collector.objects_users_num, prometheus.GaugeValue, userCount)
}
