package main

import (
	"fmt"
	"log"
	"time"

	m "github.com/grafana/grafana/pkg/models"
)

func load(c client, mail string) {

	log.Println("getting list of collectors to use")
	collectorIds := c.publicCollectorIds()
	alertCollErrors := len(collectorIds)
	// for alerting, never ask to be alerted if num coll are erroring if num is more than the actual number of collectors in the footprint
	// see also https://github.com/raintank/grafana/issues/480
	// the value used will cycle between 0 and this, so that we can see different stages of endpoints erroring
	if alertCollErrors > 10 {
		alertCollErrors = 10
	}

	for o := 1; o <= *orgs; o++ {
		u := NewUser(o, mail)
		fmt.Println(">> creating user", u.settings.Name)
		c.CreateUserForm(u.settings)
		numCollErrors := 1
		if alertCollErrors >= 2 {
			numCollErrors = (o % alertCollErrors) + 1
		}
		for e := 1; e <= 4; e++ {
			settings := m.AddEndpointCommand{
				Name: fmt.Sprintf("fake_org_%d_endpoint_%d", o, e),
				Tags: make([]string, 0),
				Monitors: []*m.AddMonitorCommand{
					pingMonitor(collectorIds, numCollErrors, e, *monHost, *email),
					dnsMonitor(collectorIds, numCollErrors, e, *monHost, *email),
					httpMonitor(collectorIds, numCollErrors, e, *monHost, *email),
				},
			}
			fmt.Println(">> creating endpoint", settings.Name)
			if err := u.c.NewEndpoint(settings); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Duration(*delay) * time.Second)
		}
	}
}
