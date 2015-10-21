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
		for e := 1; e <= 4; e++ {
			settings := m.AddEndpointCommand{
				Name: fmt.Sprintf("fake_org_%d_endpoint_%d", o, e),
				Tags: make([]string, 0),
				Monitors: []*m.AddMonitorCommand{
					pingMonitor(collectorIds, (o%10)+1, e, *email),
					dnsMonitor(collectorIds, (o%10)+1, e, *email),
					httpMonitor(collectorIds, (o%10)+1, e, *email),
				},
			}
			fmt.Println(">> creating endpoint", settings.Name)
			if err = subClient.NewEndpoint(settings); err != nil {
				log.Fatal(err)
			}
		}
	}
}
