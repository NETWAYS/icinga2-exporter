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

func TestIcinga2ApplicationCollector_Collect(t *testing.T) {
	client := &MockIcingaClient{}

	var result icinga.ApplicationResult
	data, _ := os.ReadFile("testdata/application.json")
	json.Unmarshal(data, &result)

	client.SetApplicationMetrics(result)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector := NewIcinga2ApplicationCollector(client, logger)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricCount := testutil.CollectAndCount(collector)

	expectedCount := 1
	if metricCount != expectedCount {
		t.Errorf("expected %d metrics, got %d", expectedCount, metricCount)
	}

	expectedMetrics := `
     # HELP icinga2_version_info A metric with a constant '1' value labeled by version
     # TYPE icinga2_version_info gauge
     icinga2_version_info{version="r2.15.2-1"} 1
    `

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "icinga2_version_info"); err != nil {
		t.Errorf("unexpected metric difference:\n%s", err)
	}
}
