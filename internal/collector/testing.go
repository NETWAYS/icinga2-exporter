package collector

import (
	"github.com/NETWAYS/icinga2-exporter/internal/icinga"
)

// MockIcingaClient is a simple client we use in unitests
type MockIcingaClient struct {
	perfdata    icinga.PerfdataResult
	cib         icinga.CIBResult
	api         icinga.APIResult
	application icinga.ApplicationResult
	checker     icinga.CheckerComponentResult
}

func (m *MockIcingaClient) SetPerfdataMetrics(perfdata icinga.PerfdataResult) {
	m.perfdata = perfdata
}

func (m *MockIcingaClient) GetPerfdataMetrics(endpoint string) (icinga.PerfdataResult, error) {
	return m.perfdata, nil
}

func (m *MockIcingaClient) SetCIBMetrics(cib icinga.CIBResult) {
	m.cib = cib
}

func (m *MockIcingaClient) GetCIBMetrics() (icinga.CIBResult, error) {
	return m.cib, nil
}

func (m *MockIcingaClient) SetAPIListenerMetrics(api icinga.APIResult) {
	m.api = api
}

func (m *MockIcingaClient) GetAPIListenerMetrics() (icinga.APIResult, error) {
	return m.api, nil
}

func (m *MockIcingaClient) SetCheckerComponentMetrics(checker icinga.CheckerComponentResult) {
	m.checker = checker
}

func (m *MockIcingaClient) GetCheckerComponentMetrics() (icinga.CheckerComponentResult, error) {
	return m.checker, nil
}

func (m *MockIcingaClient) SetApplicationMetrics(application icinga.ApplicationResult) {
	m.application = application
}

func (m *MockIcingaClient) GetApplicationMetrics() (icinga.ApplicationResult, error) {
	return m.application, nil
}
