package icinga

type Perfdata struct {
	IsCounter bool    `json:"counter"`
	Label     string  `json:"label"`
	Value     float64 `json:"value"`
}

type APIResult struct {
	Results []struct {
		Name     string     `json:"name"`
		Perfdata []Perfdata `json:"perfdata,omitempty"`
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
		Name     string     `json:"name"`
		Perfdata []Perfdata `json:"perfdata,omitempty"`
	} `json:"results"`
}
