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
			&m.MonitorSettingDTO{Variable: "timeout", Value: "3" },
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
			&m.MonitorSettingDTO{Variable: "timeout", Value: "3" },
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
			&m.MonitorSettingDTO{Variable: "timeout", Value: "3" },
		},
		HealthSettings: healthSettings(collectors, steps, email),
		Frequency:      10,
		Enabled:        true,
	}
}
