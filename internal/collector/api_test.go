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

func TestIcinga2APICollector_Collect(t *testing.T) {
	client := &MockIcingaClient{}

	var result icinga.APIResult
	data, _ := os.ReadFile("testdata/api.json")
	json.Unmarshal(data, &result)

	client.SetAPIListenerMetrics(result)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector := NewIcinga2APICollector(client, logger)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricCount := testutil.CollectAndCount(collector)

	expectedCount := 10
	if metricCount != expectedCount {
		t.Errorf("expected %d metrics, got %d", expectedCount, metricCount)
	}

	expectedMetrics := `
        # HELP icinga2_api_jsonrpc_sync_queue_item_rate JSON RPC sync queue item rate
        # TYPE icinga2_api_jsonrpc_sync_queue_item_rate gauge
        icinga2_api_jsonrpc_sync_queue_item_rate 3
    `

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "icinga2_api_jsonrpc_sync_queue_item_rate"); err != nil {
		t.Errorf("unexpected metric difference:\n%s", err)
	}
}
