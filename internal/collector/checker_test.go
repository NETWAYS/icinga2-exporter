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

func TestIcinga2CheckerCollector_Collect(t *testing.T) {
	client := &MockIcingaClient{}

	var result icinga.CheckerComponentResult
	data, _ := os.ReadFile("testdata/checker.json")
	json.Unmarshal(data, &result)

	client.SetCheckerComponentMetrics(result)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector := NewIcinga2CheckerCollector(client, logger)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricCount := testutil.CollectAndCount(collector)

	expectedCount := 2
	if metricCount != expectedCount {
		t.Errorf("expected %d metrics, got %d", expectedCount, metricCount)
	}

	expectedMetrics := `
        # HELP icinga2_checkercomponent_checker_pending CheckerComponent pending
        # TYPE icinga2_checkercomponent_checker_pending gauge
        icinga2_checkercomponent_checker_pending 41
    `

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "icinga2_checkercomponent_checker_pending"); err != nil {
		t.Errorf("unexpected metric difference:\n%s", err)
	}
}
