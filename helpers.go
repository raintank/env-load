package main

import m "github.com/grafana/grafana/pkg/models"

func healthSettings(collectors, steps int, email string) *m.MonitorHealthSettingDTO {
	return &m.MonitorHealthSettingDTO{
		NumCollectors: collectors,
		Steps:         steps,
		Notifications: m.MonitorNotificationSetting{
			Enabled:   true,
			Addresses: email,
		},
	}
}

func httpMonitor(collectorIds []int64, collectors, steps int, email string) *m.AddMonitorCommand {
	return &m.AddMonitorCommand{
		EndpointId:    0,
		MonitorTypeId: 1,
		CollectorIds:  collectorIds,
		CollectorTags: make([]string, 0),
		Settings: []*m.MonitorSettingDTO{
			&m.MonitorSettingDTO{Variable: "host", Value: "localhost"},
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

func pingMonitor(collectorIds []int64, collectors, steps int, email string) *m.AddMonitorCommand {
	return &m.AddMonitorCommand{
		EndpointId:    0,
		MonitorTypeId: 3,
		CollectorIds:  collectorIds,
		CollectorTags: make([]string, 0),
		Settings: []*m.MonitorSettingDTO{
			&m.MonitorSettingDTO{Variable: "hostname", Value: "localhost"},
		},
		HealthSettings: healthSettings(collectors, steps, email),
		Frequency:      10,
		Enabled:        true,
	}
}

func dnsMonitor(collectorIds []int64, collectors, steps int, email string) *m.AddMonitorCommand {
	return &m.AddMonitorCommand{
		EndpointId:    0,
		MonitorTypeId: 4,
		CollectorIds:  collectorIds,
		CollectorTags: make([]string, 0),
		Settings: []*m.MonitorSettingDTO{
			&m.MonitorSettingDTO{Variable: "name", Value: "localhost"},
			&m.MonitorSettingDTO{Variable: "type", Value: "A"},
			&m.MonitorSettingDTO{Variable: "server", Value: "8.8.8.8"},
			&m.MonitorSettingDTO{Variable: "port", Value: "53"},
			&m.MonitorSettingDTO{Variable: "protocol", Value: "udp"},
		},
		HealthSettings: healthSettings(collectors, steps, email),
		Frequency:      60,
		Enabled:        true,
	}
}
