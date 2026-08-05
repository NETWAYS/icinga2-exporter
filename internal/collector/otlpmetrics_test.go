package collector

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/NETWAYS/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestIcinga2OTLPMetricsWriterCollector_Collect(t *testing.T) {
	client := &MockIcingaClient{}

	var result icinga.PerfdataResult
	data, _ := os.ReadFile("testdata/otlp.json")
	json.Unmarshal(data, &result)

	client.SetPerfdataMetrics(result)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector := NewIcinga2OTLPMetricsCollector(client, logger)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricCount := testutil.CollectAndCount(collector)

	expectedCount := 8
	if metricCount != expectedCount {
		t.Errorf("expected %d metrics, got %d", expectedCount, metricCount)
	}

	expectedMetrics := `
  # HELP icinga2_otlpmetricswriter_work_queue_items OTLPMetricsWriter icinga2_otlpmetricswriter_work_queue_items
  # TYPE icinga2_otlpmetricswriter_work_queue_items gauge
  icinga2_otlpmetricswriter_work_queue_items{writer="otlp-metrics"} 0
  icinga2_otlpmetricswriter_work_queue_items{writer="otlp-elastic"} 0
	  `

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "icinga2_otlpmetricswriter_work_queue_items"); err != nil {
		t.Errorf("unexpected metric difference:\n%s", err)
	}
}
