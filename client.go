package main

// easy to use wrapper client around gapi.Client that provides functions specific to the env-load context.
// doesn't return errors, just fails in place when appropriate.

import (
	"log"
	"strconv"
	"strings"

	"github.com/grafana/grafana-api-golang-client"
	"github.com/grafana/grafana/pkg/api/dtos"
	m "github.com/grafana/grafana/pkg/models"
)

type client struct {
	*gapi.Client
}

func (c client) publicCollectorIds() []int64 {
	collectors, err := c.Collectors(m.GetCollectorsQuery{Public: "true"})
	if err != nil {
		log.Fatal(err)
	}
	collectorIds := make([]int64, 0)
	for _, coll := range collectors {
		collectorIds = append(collectorIds, int64(coll.Id))
	}
	return collectorIds
}

func (c client) users() []user {
	users, err := c.Client.Users()
	if err != nil {
		log.Fatal(err)
	}
	ret := make([]user, 0)
	for _, usr := range users {
		if strings.HasPrefix(usr.Name, "fake_user") {
			ret = append(ret, user{usr, nil})
		}
	}
	return ret
}

func (c client) orgs() []gapi.Org {
	orgs, err := c.Orgs()
	if err != nil {
		log.Fatal(err)
	}
	ret := make([]gapi.Org, 0)
	for _, org := range orgs {
		if strings.HasPrefix(org.Name, "fake_user") {
			ret = append(ret, org)
		}
	}
	return ret
}

func (c client) createUserForm(settings dtos.AdminCreateUserForm) {
	err := c.CreateUserForm(settings)
	if err != nil {
		log.Fatal(err)
	}
}
