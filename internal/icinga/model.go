package icinga

type Perfdata struct {
	IsCounter bool    `json:"counter"`
	Label     string  `json:"label"`
	Value     float64 `json:"value"`
}

type PerfdataResult struct {
	Results []struct {
		Name     string                                   `json:"name"`
		Perfdata []Perfdata                               `json:"perfdata,omitempty"`
		Status   map[string]map[string]map[string]float64 `json:"status"`
	} `json:"results"`
}

type CIBResult struct {
	Results []struct {
		Name   string             `json:"name"`
		Status map[string]float64 `json:"status,omitempty"`
	} `json:"results"`
}

type ApplicationResult struct {
	Results []struct {
		Name   string `json:"name"`
		Status struct {
			IcingaApplication IcingaApplication `json:"icingaapplication"`
		} `json:"status"`
	} `json:"results"`
}

type IcingaApplication struct {
	App App `json:"app"`
}

type App struct {
	EnableEventHandlers bool   `json:"enable_event_handlers"`
	EnableFlapping      bool   `json:"enable_flapping"`
	EnableHostChecks    bool   `json:"enable_host_checks"`
	EnableNotifications bool   `json:"enable_notifications"`
	EnablePerfdata      bool   `json:"enable_perfdata"`
	EnableServiceChecks bool   `json:"enable_service_checks"`
	Version             string `json:"version"`
}

type CheckerComponentResult struct {
	Results []struct {
		Name   string `json:"name"`
		Status struct {
			CheckerComponent CheckerComponent `json:"checkercomponent"`
		} `json:"status"`
	} `json:"results"`
}

type CheckerComponent struct {
	Checker struct {
		Idle    float64 `json:"idle"`
		Pending float64 `json:"pending"`
	} `json:"checker"`
}

type APIResult struct {
	Results []struct {
		Name   string `json:"name"`
		Status struct {
			API API `json:"api"`
		} `json:"status"`
	} `json:"results"`
}

type API struct {
	NumConnEndpoints    float64 `json:"num_conn_endpoints"`
	NumNotConnEndpoints float64 `json:"num_not_conn_endpoints"`
	NumEndpoints        float64 `json:"num_endpoints"`
	HTTP                struct {
		Clients float64 `json:"clients"`
	} `json:"http"`
	JSONRPC JSONRPC `json:"json_rpc"`
}

type JSONRPC struct {
	AnonymousClients   float64 `json:"anonymous_clients"`
	RelayQueueItemRate float64 `json:"relay_queue_item_rate"`
	RelayQueueItems    float64 `json:"relay_queue_items"`
	SyncQueueItemRate  float64 `json:"sync_queue_item_rate"`
	SyncQueueItems     float64 `json:"sync_queue_items"`
	WorkQueueItemRate  float64 `json:"work_queue_item_rate"`
}
