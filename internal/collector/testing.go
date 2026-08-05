package collector

import (
	"github.com/NETWAYS/icinga2-exporter/internal/icinga"
)

// MockIcingaClient is a simple client we use in unitests
type MockIcingaClient struct {
	perfdata    icinga.PerfdataResult
	cib         icinga.CIBResult
	application icinga.ApplicationResult
}

func (m *MockIcingaClient) SetPerfdataMetrics(perfdata icinga.PerfdataResult) {
	m.perfdata = perfdata
}

func (m *MockIcingaClient) GetPerfdataMetrics(endpoint string) ([]icinga.Perfdata, error) {
	return m.perfdata.Results[0].Perfdata, nil
}

func (m *MockIcingaClient) SetCIBMetrics(cib icinga.CIBResult) {
	m.cib = cib
}

func (m *MockIcingaClient) GetCIBMetrics() (icinga.CIBResult, error) {
	return m.cib, nil
}

func (m *MockIcingaClient) SetApplicationMetrics(application icinga.ApplicationResult) {
	m.application = application
}

func (m *MockIcingaClient) GetApplicationMetrics() (icinga.ApplicationResult, error) {
	return m.application, nil
}
