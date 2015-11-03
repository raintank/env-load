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

	for o := 1; o <= *orgs; o++ {
		u := NewUser(o, mail)
		fmt.Println(">> creating user", u.settings.Name)
		c.CreateUserForm(u.settings)
		for e := 1; e <= 4; e++ {
			end := getEndpointCommand(o, e, len(collectorIds), *monHost, *email)
			fmt.Println(">> creating endpoint", settings.Name)
			if err := u.c.NewEndpoint(settings); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Duration(*delay) * time.Second)
		}
	}
}
