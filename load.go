package main

import (
	"fmt"
	"log"

	"github.com/grafana/grafana-api-golang-client"
	"github.com/grafana/grafana/pkg/api/dtos"
	m "github.com/grafana/grafana/pkg/models"
)

func load(client *gapi.Client, mail string) {

	log.Println("getting list of collectors to use")
	collectors, err := client.Collectors(m.GetCollectorsQuery{Public: "true"})
	if err != nil {
		log.Fatal(err)
	}
	collectorIds := make([]int64, 0)
	for _, coll := range collectors {
		collectorIds = append(collectorIds, int64(coll.Id))
	}
	alertCollErrors := len(collectorIds)
	// for alerting, never ask to be alerted if num coll are erroring if num is more than the actual number of collectors in the footprint
	// see also https://github.com/raintank/grafana/issues/480
	// the value used will cycle between 0 and this, so that we can see different stages of endpoints erroring
	if alertCollErrors > 10 {
		alertCollErrors = 10
	}

	for o := 1; o <= *orgs; o++ {
		user := fmt.Sprintf("fake_user_%d", o)
		pass := fmt.Sprintf("fake_pass_%d", o)
		settings := dtos.AdminCreateUserForm{
			Email:    mail,
			Login:    user,
			Name:     user,
			Password: pass,
		}
		fmt.Println(">> creating user", user)
		err = client.CreateUserForm(settings)
		if err != nil {
			log.Fatal(err)
		}
		subClient, err := gapi.New(fmt.Sprintf("%s:%s", user, pass), *host)
		if err != nil {
			log.Fatal(err)
		}
		numCollErrors := 1
		if alertCollErrors >= 2 {
			numCollErrors = (o % alertCollErrors) + 1
		}
		for e := 1; e <= 4; e++ {
			settings := m.AddEndpointCommand{
				Name: fmt.Sprintf("fake_org_%d_endpoint_%d", o, e),
				Tags: make([]string, 0),
				Monitors: []*m.AddMonitorCommand{
					pingMonitor(collectorIds, numCollErrors, e, *email),
					dnsMonitor(collectorIds, numCollErrors, e, *email),
					httpMonitor(collectorIds, numCollErrors, e, *email),
				},
			}
			fmt.Println(">> creating endpoint", settings.Name)
			if err = subClient.NewEndpoint(settings); err != nil {
				log.Fatal(err)
			}
		}
	}
}
