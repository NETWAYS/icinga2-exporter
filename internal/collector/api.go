package collector

import (
	"log/slog"

	"github.com/NETWAYS/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2APICollector struct {
	icingaClient                      icinga.IcingaClient
	logger                            *slog.Logger
	api_num_conn_endpoints            *prometheus.Desc
	api_num_not_conn_endpoints        *prometheus.Desc
	api_num_endpoints                 *prometheus.Desc
	api_num_http_clients              *prometheus.Desc
	api_jsonrpc_anonymous_clients     *prometheus.Desc
	api_jsonrpc_relay_queue_item_rate *prometheus.Desc
	api_jsonrpc_relay_queue_items     *prometheus.Desc
	api_jsonrpc_sync_queue_item_rate  *prometheus.Desc
	api_jsonrpc_sync_queue_items      *prometheus.Desc
	api_jsonrpc_work_queue_item_rate  *prometheus.Desc
}

func NewIcinga2APICollector(client icinga.IcingaClient, logger *slog.Logger) *Icinga2APICollector {
	return &Icinga2APICollector{
		icingaClient:               client,
		logger:                     logger,
		api_num_conn_endpoints:     prometheus.NewDesc("icinga2_api_num_conn_endpoints", "Number of connected endpoints", nil, nil),
		api_num_endpoints:          prometheus.NewDesc("icinga2_api_num_endpoints", "Number of endpoints", nil, nil),
		api_num_not_conn_endpoints: prometheus.NewDesc("icinga2_api_num_not_conn_endpoints", "Number of not connected endpoints", nil, nil),
		api_num_http_clients:       prometheus.NewDesc("icinga2_api_num_http_clients", "Number of HTTP clients", nil, nil),

		api_jsonrpc_anonymous_clients:     prometheus.NewDesc("icinga2_api_jsonrpc_anonymous_clients", "JSON RPC anonymous clients", nil, nil),
		api_jsonrpc_relay_queue_item_rate: prometheus.NewDesc("icinga2_api_jsonrpc_relay_queue_item_rate", "JSON RPC relay queue item rate", nil, nil),
		api_jsonrpc_relay_queue_items:     prometheus.NewDesc("icinga2_api_jsonrpc_relay_queue_items", "JSON RPC relay queue items", nil, nil),
		api_jsonrpc_sync_queue_item_rate:  prometheus.NewDesc("icinga2_api_jsonrpc_sync_queue_item_rate", "JSON RPC sync queue item rate", nil, nil),
		api_jsonrpc_sync_queue_items:      prometheus.NewDesc("icinga2_api_jsonrpc_sync_queue_items", "JSON RPC sync queue items", nil, nil),
		api_jsonrpc_work_queue_item_rate:  prometheus.NewDesc("icinga2_api_jsonrpc_work_queue_item_rate", "JSON RPC work queue item rate", nil, nil),
	}
}

func (collector *Icinga2APICollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.api_num_conn_endpoints

	ch <- collector.api_num_not_conn_endpoints

	ch <- collector.api_num_endpoints

	ch <- collector.api_num_http_clients

	ch <- collector.api_jsonrpc_anonymous_clients

	ch <- collector.api_jsonrpc_relay_queue_item_rate

	ch <- collector.api_jsonrpc_relay_queue_items

	ch <- collector.api_jsonrpc_sync_queue_item_rate

	ch <- collector.api_jsonrpc_sync_queue_items

	ch <- collector.api_jsonrpc_work_queue_item_rate
}

func (collector *Icinga2APICollector) Collect(ch chan<- prometheus.Metric) {
	result, err := collector.icingaClient.GetAPIListenerMetrics()

	if err != nil {
		collector.logger.Error("Could not retrieve ApiListener metrics", "error", err.Error())
		return
	}

	if len(result.Results) < 1 {
		collector.logger.Debug("No results for ApiListener metrics")
		return
	}

	r := result.Results[0]

	ch <- prometheus.MustNewConstMetric(collector.api_num_conn_endpoints, prometheus.GaugeValue, r.Status.API.NumConnEndpoints)

	ch <- prometheus.MustNewConstMetric(collector.api_num_not_conn_endpoints, prometheus.GaugeValue, r.Status.API.NumNotConnEndpoints)

	ch <- prometheus.MustNewConstMetric(collector.api_num_endpoints, prometheus.GaugeValue, r.Status.API.NumEndpoints)

	ch <- prometheus.MustNewConstMetric(collector.api_num_http_clients, prometheus.GaugeValue, r.Status.API.HTTP.Clients)

	ch <- prometheus.MustNewConstMetric(collector.api_jsonrpc_anonymous_clients, prometheus.GaugeValue, r.Status.API.JSONRPC.AnonymousClients)

	ch <- prometheus.MustNewConstMetric(collector.api_jsonrpc_relay_queue_item_rate, prometheus.GaugeValue, r.Status.API.JSONRPC.SyncQueueItemRate)

	ch <- prometheus.MustNewConstMetric(collector.api_jsonrpc_relay_queue_items, prometheus.GaugeValue, r.Status.API.JSONRPC.RelayQueueItems)

	ch <- prometheus.MustNewConstMetric(collector.api_jsonrpc_sync_queue_item_rate, prometheus.GaugeValue, r.Status.API.JSONRPC.SyncQueueItemRate)

	ch <- prometheus.MustNewConstMetric(collector.api_jsonrpc_sync_queue_items, prometheus.GaugeValue, r.Status.API.JSONRPC.SyncQueueItems)

	ch <- prometheus.MustNewConstMetric(collector.api_jsonrpc_work_queue_item_rate, prometheus.GaugeValue, r.Status.API.JSONRPC.WorkQueueItemRate)
}
