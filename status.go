package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grafana/grafana-api-golang-client"
	"github.com/grafana/grafana/pkg/models"
)

// based on clean, but obviously without cleaning anything
func status(client *gapi.Client) {

	log.Println("currently existing fake users and their endpoints")
	users, err := client.Users()
	if err != nil {
		log.Fatal(err)
	}
	for _, usr := range users {
		if strings.HasPrefix(usr.Name, "fake_user") {
			fmt.Println(usr.Id, usr.Name)
			subClient, err := gapi.New(fmt.Sprintf("%s:%s", usr.Name, strings.Replace(usr.Name, "user", "pass", 1)), *host)
			if err != nil {
				log.Fatal(err)
			}
			settings := models.GetEndpointsQuery{}
			endpoints, err := subClient.Endpoints(settings)
			if err != nil {
				log.Fatal(err)
			}
			for _, e := range endpoints {
				fmt.Println(e.Id, e.Name, e.OrgId)
			}
		}
	}

	log.Println("currently existing fake orgs")
	orgs, err := client.Orgs()
	if err != nil {
		log.Fatal(err)
	}
	for _, org := range orgs {
		if strings.HasPrefix(org.Name, "fake_user") {
			fmt.Println(org.Id, org.Name)
		}
	}
}
