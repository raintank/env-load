package main

import (
	"fmt"
	"log"
)

// deleting endpoints is finicky.
// we can't just get all endpoints as admin like we can for users/orgs,
// and while grafana lets us get endpoints with an org-id query, it overrides
// the org-id to the active org of the user.
// so we must login as that user and rely on the active org of the user, so
// don't change that (by logging in as that user and changing the org)!
func clean(c client) {

	log.Println("removing fake users and their endpoints...")
	for _, usr := range c.users() {
		fmt.Println(usr.Id, usr.Name)
		for _, e := range usr.endpoints() {
			fmt.Println(e.Id, e.OrgId, e.Name)
			usr.deleteEndpoint(e.Id)
		}
		c.DeleteUser(usr.Id)
	}

	log.Println("removing fake orgs...")
	for _, org := range c.orgs() {
		fmt.Println(org.Id, org.Name)
		err := c.DeleteOrg(org.Id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
