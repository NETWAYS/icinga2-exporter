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
	icingaTestDataAPP1 = "testdata/app1.json"
	icingaTestDataAPI1 = "testdata/api1.json"
	// icingaTestDataCheck1 = "testdata/checker1.json"
	icingaTestDataPerf1 = "testdata/otlp.json"
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

func Test_GetApplicationMetrics(t *testing.T) {
	testcases := map[string]struct {
		expected ApplicationResult
		server   *httptest.Server
	}{
		"application": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(loadTestdata(icingaTestDataAPP1))
			})),
			expected: ApplicationResult{
				Results: []struct {
					Name   string `json:"name"`
					Status struct {
						IcingaApplication IcingaApplication `json:"icingaapplication"`
					} `json:"status"`
				}{
					{
						Name: "IcingaApplication",
						Status: struct {
							IcingaApplication IcingaApplication `json:"icingaapplication"`
						}{
							IcingaApplication: IcingaApplication{
								App: App{
									EnableEventHandlers: true,
									Version:             "r2.15.2-1",
								},
							},
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

			actual, err := cli.GetApplicationMetrics()

			if err != nil {
				t.Fatalf("did not expect error got:\n %+v", err)
			}

			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatalf("expected:\n %+v \ngot:\n %+v", test.expected, actual)
			}
		})
	}
}

func Test_GetPerfdataMetrics(t *testing.T) {
	testcases := map[string]struct {
		expected PerfdataResult
		endpoint string
		server   *httptest.Server
	}{
		"otlp": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(loadTestdata(icingaTestDataPerf1))
			})),
			endpoint: EndpointOTLPMetricsWriter,
			expected: PerfdataResult{
				Results: []struct {
					Name     string                                   `json:"name"`
					Perfdata []Perfdata                               `json:"perfdata,omitempty"`
					Status   map[string]map[string]map[string]float64 `json:"status"`
				}{
					{
						Name: "OTLPMetricsWriter",
						Perfdata: []Perfdata{
							{
								IsCounter: false,
								Label:     "otlpmetricswriter_otlp-elastic_work_queue_items",
								Value:     0,
							},
							{
								IsCounter: false,
								Label:     "otlpmetricswriter_otlp-elastic_work_queue_item_rate",
								Value:     0.5833333333333334,
							},
						},
						Status: map[string]map[string]map[string]float64{
							"otlpmetricswriter": {
								"otlp-elastic": {
									"data_buffer_items": 48095,
								},
								"otlp-metrics": {
									"data_buffer_items": 153,
								},
							},
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

			actual, err := cli.GetPerfdataMetrics(test.endpoint)

			if err != nil {
				t.Fatalf("did not expect error got:\n %+v", err)
			}

			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatalf("expected:\n %+v \ngot:\n %+v", test.expected, actual)
			}
		})
	}
}
