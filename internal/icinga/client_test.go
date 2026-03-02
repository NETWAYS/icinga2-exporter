package icinga

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	icingaTestDataCIB1 = "testdata/cib1.json"
)

func loadTestdata(filepath string) []byte {
	data, _ := os.ReadFile(filepath)
	return data
}

func testConfig(ts *httptest.Server) Config {
	u, _ := url.Parse(ts.URL)
	return Config{
		BasicAuthUsername: "",
		BasicAuthPassword: "",
		CAFile:            "",
		CertFile:          "",
		KeyFile:           "",
		Insecure:          true,
		CacheTTL:          0 * time.Second,
		IcingaAPIURI:      *u,
	}
}

func Test_GetCIBMetrics(t *testing.T) {
	testcases := map[string]struct {
		expected CIBResult
		server   *httptest.Server
	}{
		"cib": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(loadTestdata(icingaTestDataCIB1))
			})),
			expected: CIBResult{
				Results: []struct {
					Name   string             `json:"name"`
					Status map[string]float64 `json:"status,omitempty"`
				}{
					{
						Name: "CIB",
						Status: map[string]float64{
							"active_host_checks": 0.123,
							"uptime":             123.456,
						},
					},
				},
			},
		},
	}

	for name, test := range testcases {
		t.Run(name, func(t *testing.T) {
			defer test.server.Close()

			cfg := testConfig(test.server)

			cli, _ := NewClient(cfg)

			actual, err := cli.GetCIBMetrics()

			if err != nil {
				t.Fatalf("did not expect error got:\n %+v", err)
			}

			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatalf("expected:\n %+v \ngot:\n %+v", test.expected, actual)
			}
		})
	}
}
