package main

import m "github.com/grafana/grafana/pkg/models"

func healthSettings(collectors, steps int, email string) *m.MonitorHealthSettingDTO {
	hs := &m.MonitorHealthSettingDTO{
		NumCollectors: collectors,
		Steps:         steps,
		Notifications: m.MonitorNotificationSetting{
			Enabled:   false,
			Addresses: email,
		},
	}
	if email != "" {
		hs.Notifications.Enabled = true
	}
	return hs
}

func httpMonitor(collectorIds []int64, collectors, steps int, host, email string) *m.AddMonitorCommand {
	return &m.AddMonitorCommand{
		EndpointId:    1,
		MonitorTypeId: 1,
		CollectorIds:  collectorIds,
		CollectorTags: make([]string, 0),
		Settings: []*m.MonitorSettingDTO{
			&m.MonitorSettingDTO{Variable: "host", Value: host},
			&m.MonitorSettingDTO{Variable: "port", Value: "80"},
			&m.MonitorSettingDTO{Variable: "path", Value: "/"},
			&m.MonitorSettingDTO{Variable: "method", Value: "GET"},
			&m.MonitorSettingDTO{Variable: "expectRegex", Value: ""},
			&m.MonitorSettingDTO{Variable: "headers", Value: "Accept-Encoding: gzip\nUser-Agent: raintank collector\n"},
		},
		HealthSettings: healthSettings(collectors, steps, email),
		Frequency:      10,
		Enabled:        true,
	}
}

func pingMonitor(collectorIds []int64, collectors, steps int, host, email string) *m.AddMonitorCommand {
	return &m.AddMonitorCommand{
		EndpointId:    1,
		MonitorTypeId: 3,
		CollectorIds:  collectorIds,
		CollectorTags: make([]string, 0),
		Settings: []*m.MonitorSettingDTO{
			&m.MonitorSettingDTO{Variable: "hostname", Value: host},
		},
		HealthSettings: healthSettings(collectors, steps, email),
		Frequency:      10,
		Enabled:        true,
	}
}

func dnsMonitor(collectorIds []int64, collectors, steps int, host, email string) *m.AddMonitorCommand {
	return &m.AddMonitorCommand{
		EndpointId:    1,
		MonitorTypeId: 4,
		CollectorIds:  collectorIds,
		CollectorTags: make([]string, 0),
		Settings: []*m.MonitorSettingDTO{
			&m.MonitorSettingDTO{Variable: "name", Value: host},
			&m.MonitorSettingDTO{Variable: "type", Value: "A"},
			&m.MonitorSettingDTO{Variable: "server", Value: "8.8.8.8"},
			&m.MonitorSettingDTO{Variable: "port", Value: "53"},
			&m.MonitorSettingDTO{Variable: "protocol", Value: "udp"},
		},
		HealthSettings: healthSettings(collectors, steps, email),
		Frequency:      10,
		Enabled:        true,
	}
}
func endpointCommand(orgI, endpI, numCollectors int, host, email string) m.AddEndpointCommand {
	alertCollErrors = numCollectors
	// for alerting, never ask to be alerted if num coll are erroring if num is more than the actual number of collectors in the footprint
	// see also https://github.com/raintank/grafana/issues/480
	// the value used will cycle between 0 and this, so that we can see different stages of endpoints erroring
	if alertCollErrors > 10 {
		alertCollErrors = 10
	}
	numCollErrors := 1
	if alertCollErrors >= 2 {
		numCollErrors = (o % alertCollErrors) + 1
	}
	end := getEndpointCommand(o, e, *monHost, *email)
	return m.AddEndpointCommand{
		Name: fmt.Sprintf("fake_org_%d_endpoint_%d", o, e),
		Tags: make([]string, 0),
		Monitors: []*m.AddMonitorCommand{
			pingMonitor(collectorIds, numCollErrors, e, host, email),
			dnsMonitor(collectorIds, numCollErrors, e, host, email),
			httpMonitor(collectorIds, numCollErrors, e, host, email),
		},
	}
}
