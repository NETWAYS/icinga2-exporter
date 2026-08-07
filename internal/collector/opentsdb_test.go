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

func TestIcinga2OpenTSDBWriterCollector_Collect(t *testing.T) {
	client := &MockIcingaClient{}

	var result icinga.PerfdataResult
	data, _ := os.ReadFile("testdata/opentsdb.json")
	json.Unmarshal(data, &result)

	client.SetPerfdataMetrics(result)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector := NewIcinga2OpenTSDBCollector(client, logger)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricCount := testutil.CollectAndCount(collector)

	expectedCount := 2
	if metricCount != expectedCount {
		t.Errorf("expected %d metrics, got %d", expectedCount, metricCount)
	}

	expectedMetrics := `
  # HELP icinga2_opentsdbwriter_work_queue_items OpenTSDBWriter icinga2_opentsdbwriter_work_queue_items
  # TYPE icinga2_opentsdbwriter_work_queue_items gauge
  icinga2_opentsdbwriter_work_queue_items{writer="opentsdb"} 4
	  `

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "icinga2_opentsdbwriter_work_queue_items"); err != nil {
		t.Errorf("unexpected metric difference:\n%s", err)
	}
}
