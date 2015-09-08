package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grafana/grafana-api-golang-client"
	"github.com/grafana/grafana/pkg/models"
)

// deleting endpoints is finicky.
// we can't just get all endpoints as admin like we can for users/orgs,
// and while grafana lets us get endpoints with an org-id query, it overrides
// the org-id to the active org of the user.
// so we must login as that user and rely on the active org of the user, so
// don't change that (by logging in as that user and changing the org)!
func clean(client *gapi.Client) {

	log.Println("removing fake users and their endpoints...")
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
			// no point in setting org-id in settings, because grafana overrides it anyway
			settings := models.GetEndpointsQuery{}
			endpoints, err := subClient.Endpoints(settings)
			if err != nil {
				log.Fatal(err)
			}

			for _, e := range endpoints {
				fmt.Println(e.Name, e.OrgId, e.Id)
				err = subClient.DeleteEndpoint(e.Id)
				if err != nil {
					log.Fatal(err)
				}
			}

			client.DeleteUser(usr.Id)
		}
	}

	log.Println("removing fake orgs...")
	orgs, err := client.Orgs()
	if err != nil {
		log.Fatal(err)
	}
	for _, org := range orgs {
		if strings.HasPrefix(org.Name, "fake_user") {
			fmt.Println(org.Id, org.Name)
			err = client.DeleteOrg(org.Id)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
