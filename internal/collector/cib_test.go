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

func TestIcinga2CIBCollector_Collect(t *testing.T) {
	client := &MockIcingaClient{}

	var result icinga.CIBResult
	data, _ := os.ReadFile("testdata/cib.json")
	json.Unmarshal(data, &result)

	client.SetCIBMetrics(result)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector := NewIcinga2CIBCollector(client, logger)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricCount := testutil.CollectAndCount(collector)

	expectedCount := 46
	if metricCount != expectedCount {
		t.Errorf("expected %d metrics, got %d", expectedCount, metricCount)
	}

	expectedMetrics := `
        # HELP icinga2_active_service_checks_5min Active service checks last 5min
        # TYPE icinga2_active_service_checks_5min gauge
        icinga2_active_service_checks_5min 10
    `

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "icinga2_active_service_checks_5min"); err != nil {
		t.Errorf("unexpected metric difference:\n%s", err)
	}
}
