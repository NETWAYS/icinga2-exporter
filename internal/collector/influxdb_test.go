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

func TestIcinga2InfluxDBWriterCollector_Collect(t *testing.T) {
	client := &MockIcingaClient{}

	var result icinga.PerfdataResult
	data, _ := os.ReadFile("testdata/influx.json")
	json.Unmarshal(data, &result)

	client.SetPerfdataMetrics(result)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector := NewIcinga2InfluxDBCollector(client, logger)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricCount := testutil.CollectAndCount(collector)

	expectedCount := 3
	if metricCount != expectedCount {
		t.Errorf("expected %d metrics, got %d", expectedCount, metricCount)
	}

	expectedMetrics := `
  # HELP icinga2_influxdbwriter_data_buffer_items InfluxDBWriter icinga2_influxdbwriter_data_buffer_items
  # TYPE icinga2_influxdbwriter_data_buffer_items gauge
  icinga2_influxdbwriter_data_buffer_items{writer="influxdb"} 74
	  `

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "icinga2_influxdbwriter_data_buffer_items"); err != nil {
		t.Errorf("unexpected metric difference:\n%s", err)
	}
}
