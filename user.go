package main

// easy to use wrapper user around user that provides functions specific to the env-load context.
// doesn't return errors, just fails in place when appropriate.

import (
	"fmt"
	"log"
	"strings"

	"github.com/grafana/grafana-api-golang-client"
	"github.com/grafana/grafana/pkg/api/dtos"
	m "github.com/grafana/grafana/pkg/models"
)

type user struct {
	gapi.User
	c *client
}

func (u *user) assureClient() {
	if u.c == nil {
		subClient, err := gapi.New(fmt.Sprintf("%s:%s", u.Name, strings.Replace(u.Name, "user", "pass", 1)), *host)
		if err != nil {
			log.Fatal(err)
		}
		u.c = &client{subClient}
	}
}

func (u user) endpoints() []gapi.Endpoint {
	u.assureClient()
	// no point in setting org-id in settings, because grafana overrides it anyway
	settings := m.GetEndpointsQuery{}
	endpoints, err := u.c.Endpoints(settings)
	if err != nil {
		log.Fatal(err)
	}
	return endpoints
}

func (u user) deleteEndpoint(id int64) {
	u.assureClient()
	err := u.c.DeleteEndpoint(id)
	if err != nil {
		log.Fatal(err)
	}
}

type newUser struct {
	settings dtos.AdminCreateUserForm
	c        *client
}

func NewUser(num int, mail string) newUser {
	user := fmt.Sprintf("fake_user_%d", num)
	pass := fmt.Sprintf("fake_pass_%d", num)
	settings := dtos.AdminCreateUserForm{
		Email:    mail,
		Login:    user,
		Name:     user,
		Password: pass,
	}
	c, err := gapi.New(fmt.Sprintf("%s:%s", user, pass), *host)
	if err != nil {
		log.Fatal(err)
	}
	return newUser{settings, &client{c}}
}
