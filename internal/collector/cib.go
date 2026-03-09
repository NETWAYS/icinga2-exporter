package collector

import (
	"log/slog"

	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
)

type Icinga2CIBCollector struct {
	icingaClient *icinga.Client
	logger       *slog.Logger

	// Icinga Statistics
	uptime                    *prometheus.Desc
	avg_execution_time        *prometheus.Desc
	avg_latency               *prometheus.Desc
	max_execution_time        *prometheus.Desc
	max_latency               *prometheus.Desc
	min_execution_time        *prometheus.Desc
	min_latency               *prometheus.Desc
	current_concurrent_checks *prometheus.Desc
	current_pending_callbacks *prometheus.Desc
	remote_check_queue        *prometheus.Desc

	// Active Checks
	active_host_checks          *prometheus.Desc
	active_host_checks_15min    *prometheus.Desc
	active_host_checks_1min     *prometheus.Desc
	active_host_checks_5min     *prometheus.Desc
	active_service_checks       *prometheus.Desc
	active_service_checks_15min *prometheus.Desc
	active_service_checks_1min  *prometheus.Desc
	active_service_checks_5min  *prometheus.Desc
	// Passive Checks
	passive_host_checks          *prometheus.Desc
	passive_host_checks_15min    *prometheus.Desc
	passive_host_checks_1min     *prometheus.Desc
	passive_host_checks_5min     *prometheus.Desc
	passive_service_checks       *prometheus.Desc
	passive_service_checks_15min *prometheus.Desc
	passive_service_checks_1min  *prometheus.Desc
	passive_service_checks_5min  *prometheus.Desc

	// Num Hosts
	num_hosts_up           *prometheus.Desc
	num_hosts_down         *prometheus.Desc
	num_hosts_acknowledged *prometheus.Desc
	num_hosts_flapping     *prometheus.Desc
	num_hosts_handled      *prometheus.Desc
	num_hosts_in_downtime  *prometheus.Desc
	num_hosts_pending      *prometheus.Desc
	num_hosts_problem      *prometheus.Desc
	num_hosts_unreachable  *prometheus.Desc

	// Num Services
	num_services_ok           *prometheus.Desc
	num_services_critical     *prometheus.Desc
	num_services_acknowledged *prometheus.Desc
	num_services_flapping     *prometheus.Desc
	num_services_handled      *prometheus.Desc
	num_services_in_downtime  *prometheus.Desc
	num_services_pending      *prometheus.Desc
	num_services_problem      *prometheus.Desc
	num_services_unknown      *prometheus.Desc
	num_services_unreachable  *prometheus.Desc
	num_services_warning      *prometheus.Desc
}

func NewIcinga2CIBCollector(client *icinga.Client, logger *slog.Logger) *Icinga2CIBCollector {
	return &Icinga2CIBCollector{
		icingaClient: client,
		logger:       logger,

		// Icinga Statistics
		uptime:                    prometheus.NewDesc("icinga2_uptime", "Uptime of the instance", nil, nil),
		avg_execution_time:        prometheus.NewDesc("icinga2_avg_execution_time", "Average execution time", nil, nil),
		avg_latency:               prometheus.NewDesc("icinga2_avg_latency", "Average latency", nil, nil),
		max_execution_time:        prometheus.NewDesc("icinga2_max_execution_time", "Maximum execution time", nil, nil),
		max_latency:               prometheus.NewDesc("icinga2_max_latency", "Maximum latency", nil, nil),
		min_execution_time:        prometheus.NewDesc("icinga2_min_execution_time", "Minimum execution time", nil, nil),
		min_latency:               prometheus.NewDesc("icinga2_min_latency", "Minimum latency", nil, nil),
		current_concurrent_checks: prometheus.NewDesc("icinga2_current_concurrent_checks", "Current concurrent checks", nil, nil),
		current_pending_callbacks: prometheus.NewDesc("icinga2_current_pending_callbacks", "Current pending callbacks", nil, nil),
		remote_check_queue:        prometheus.NewDesc("icinga2_remote_check_queue", "Remote check queue size", nil, nil),

		// Active Checks
		active_host_checks:          prometheus.NewDesc("icinga2_active_host_checks", "Active host checks", nil, nil),
		active_host_checks_15min:    prometheus.NewDesc("icinga2_active_host_checks_15min", "Active host checks last 15min", nil, nil),
		active_host_checks_1min:     prometheus.NewDesc("icinga2_active_host_checks_1min", "Active host checks last 1min", nil, nil),
		active_host_checks_5min:     prometheus.NewDesc("icinga2_active_host_checks_5min", "Active host checks last 5min", nil, nil),
		active_service_checks:       prometheus.NewDesc("icinga2_active_service_checks", "Active service checks", nil, nil),
		active_service_checks_15min: prometheus.NewDesc("icinga2_active_service_checks_15min", "Active service checks last 15min", nil, nil),
		active_service_checks_1min:  prometheus.NewDesc("icinga2_active_service_checks_1min", "Active service checks last 1min", nil, nil),
		active_service_checks_5min:  prometheus.NewDesc("icinga2_active_service_checks_5min", "Active service checks last 5min", nil, nil),
		// Passive Checks
		passive_host_checks:          prometheus.NewDesc("icinga2_passive_host_checks", "Passive host checks", nil, nil),
		passive_host_checks_15min:    prometheus.NewDesc("icinga2_passive_host_checks_15min", "Passive host checks last 15min", nil, nil),
		passive_host_checks_1min:     prometheus.NewDesc("icinga2_passive_host_checks_1min", "Passive host checks last 1min", nil, nil),
		passive_host_checks_5min:     prometheus.NewDesc("icinga2_passive_host_checks_5min", "Passive host checks last 5min", nil, nil),
		passive_service_checks:       prometheus.NewDesc("icinga2_passive_service_checks", "Passive service checks", nil, nil),
		passive_service_checks_15min: prometheus.NewDesc("icinga2_passive_service_checks_15min", "Passive service checks last 15min", nil, nil),
		passive_service_checks_1min:  prometheus.NewDesc("icinga2_passive_service_checks_1min", "Passive service checks last 1min", nil, nil),
		passive_service_checks_5min:  prometheus.NewDesc("icinga2_passive_service_checks_5min", "Passive service checks last 5min", nil, nil),

		// Num Hosts
		num_hosts_up:           prometheus.NewDesc("icinga2_num_hosts_up", "Number of hosts Up", nil, nil),
		num_hosts_down:         prometheus.NewDesc("icinga2_num_hosts_down", "Number of hosts Down", nil, nil),
		num_hosts_acknowledged: prometheus.NewDesc("icinga2_num_hosts_acknowledged", "Number of hosts acknowledged", nil, nil),
		num_hosts_flapping:     prometheus.NewDesc("icinga2_num_hosts_flapping", "Number of hosts flapping", nil, nil),
		num_hosts_handled:      prometheus.NewDesc("icinga2_num_hosts_handled", "Number of hosts handled", nil, nil),
		num_hosts_in_downtime:  prometheus.NewDesc("icinga2_num_hosts_in_downtime", "Number of hosts in downtime", nil, nil),
		num_hosts_pending:      prometheus.NewDesc("icinga2_num_hosts_pending", "Number of hosts pending", nil, nil),
		num_hosts_problem:      prometheus.NewDesc("icinga2_num_hosts_problem", "Number of hosts with problem", nil, nil),
		num_hosts_unreachable:  prometheus.NewDesc("icinga2_num_hosts_unreachable", "Number of hosts unreachable", nil, nil),
		// Num Services
		num_services_ok:           prometheus.NewDesc("icinga2_num_services_ok", "Number of services OK", nil, nil),
		num_services_critical:     prometheus.NewDesc("icinga2_num_services_critical", "Number of services Critical", nil, nil),
		num_services_acknowledged: prometheus.NewDesc("icinga2_num_services_acknowledged", "Number of services acknowledged", nil, nil),
		num_services_flapping:     prometheus.NewDesc("icinga2_num_services_flapping", "Number of services flapping", nil, nil),
		num_services_handled:      prometheus.NewDesc("icinga2_num_services_handled", "Number of services handled", nil, nil),
		num_services_in_downtime:  prometheus.NewDesc("icinga2_num_services_in_downtime", "Number of services in downtime", nil, nil),
		num_services_pending:      prometheus.NewDesc("icinga2_num_services_pending", "Number of services pending", nil, nil),
		num_services_problem:      prometheus.NewDesc("icinga2_num_services_problem", "Number of services with problem", nil, nil),
		num_services_unknown:      prometheus.NewDesc("icinga2_num_services_unknown", "Number of services unknown", nil, nil),
		num_services_unreachable:  prometheus.NewDesc("icinga2_num_services_unreachable", "Number of services unreachable", nil, nil),
		num_services_warning:      prometheus.NewDesc("icinga2_num_services_warning", "Number of services warning", nil, nil),
	}
}

func (collector *Icinga2CIBCollector) Describe(ch chan<- *prometheus.Desc) {
	// Icinga Statistics
	ch <- collector.uptime
	ch <- collector.avg_execution_time
	ch <- collector.avg_latency
	ch <- collector.max_execution_time
	ch <- collector.max_latency
	ch <- collector.min_execution_time
	ch <- collector.min_latency
	ch <- collector.current_concurrent_checks
	ch <- collector.current_pending_callbacks
	ch <- collector.remote_check_queue

	// Active Checks
	ch <- collector.active_host_checks
	ch <- collector.active_host_checks_15min
	ch <- collector.active_host_checks_1min
	ch <- collector.active_host_checks_5min
	ch <- collector.active_service_checks
	ch <- collector.active_service_checks_15min
	ch <- collector.active_service_checks_1min
	ch <- collector.active_service_checks_5min
	// Passive Checks
	ch <- collector.passive_host_checks
	ch <- collector.passive_host_checks_15min
	ch <- collector.passive_host_checks_1min
	ch <- collector.passive_host_checks_5min
	ch <- collector.passive_service_checks
	ch <- collector.passive_service_checks_15min
	ch <- collector.passive_service_checks_1min
	ch <- collector.passive_service_checks_5min

	// Num Hosts
	ch <- collector.num_hosts_up
	ch <- collector.num_hosts_down
	ch <- collector.num_hosts_acknowledged
	ch <- collector.num_hosts_flapping
	ch <- collector.num_hosts_handled
	ch <- collector.num_hosts_in_downtime
	ch <- collector.num_hosts_pending
	ch <- collector.num_hosts_problem
	ch <- collector.num_hosts_unreachable
	// Num Services
	ch <- collector.num_services_ok
	ch <- collector.num_services_critical
	ch <- collector.num_services_acknowledged
	ch <- collector.num_services_flapping
	ch <- collector.num_services_handled
	ch <- collector.num_services_in_downtime
	ch <- collector.num_services_pending
	ch <- collector.num_services_problem
	ch <- collector.num_services_unknown
	ch <- collector.num_services_unreachable
	ch <- collector.num_services_warning
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
	if v, ok := r.Status["current_concurrent_checks"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.current_concurrent_checks, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["current_pending_callbacks"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.current_pending_callbacks, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["remote_check_queue"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.remote_check_queue, prometheus.GaugeValue, v)
	}

	// Active Checks
	if v, ok := r.Status["active_host_checks"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.active_host_checks, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["active_host_checks_15min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.active_host_checks_15min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["active_host_checks_1min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.active_host_checks_1min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["active_host_checks_5min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.active_host_checks_5min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["active_service_checks"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.active_service_checks, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["active_service_checks_15min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.active_service_checks_15min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["active_service_checks_1min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.active_service_checks_1min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["active_service_checks_5min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.active_service_checks_5min, prometheus.GaugeValue, v)
	}

	// Passive Checks
	if v, ok := r.Status["passive_host_checks"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.passive_host_checks, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["passive_host_checks_15min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.passive_host_checks_15min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["passive_host_checks_1min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.passive_host_checks_1min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["passive_host_checks_5min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.passive_host_checks_5min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["passive_service_checks"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.passive_service_checks, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["passive_service_checks_15min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.passive_service_checks_15min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["passive_service_checks_1min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.passive_service_checks_1min, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["passive_service_checks_5min"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.passive_service_checks_5min, prometheus.GaugeValue, v)
	}

	// Hosts
	if v, ok := r.Status["num_hosts_up"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_up, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_down"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_down, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_acknowledged"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_acknowledged, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_flapping"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_flapping, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_handled"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_handled, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_in_downtime"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_in_downtime, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_pending"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_pending, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_problem"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_problem, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_hosts_unreachable"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_hosts_unreachable, prometheus.GaugeValue, v)
	}

	// Services
	if v, ok := r.Status["num_services_ok"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_ok, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_critical"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_critical, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_acknowledged"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_acknowledged, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_flapping"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_flapping, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_handled"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_handled, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_in_downtime"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_in_downtime, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_pending"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_pending, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_problem"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_problem, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_unreachable"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_unreachable, prometheus.GaugeValue, v)
	}
	if v, ok := r.Status["num_services_warning"]; ok {
		ch <- prometheus.MustNewConstMetric(collector.num_services_warning, prometheus.GaugeValue, v)
	}
}
