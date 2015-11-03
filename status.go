package main

import (
	"fmt"
	"log"
)

// based on clean, but obviously without cleaning anything
func status(c client) {

	log.Println("currently existing fake users and their endpoints")
	for _, usr := range c.users() {
		fmt.Println(usr.Id, usr.Name)
		for _, e := range usr.endpoints() {
			fmt.Println(e.Id, e.OrgId, e.Name)
		}
	}

	log.Println("currently existing fake orgs")
	for _, org := range c.orgs() {
		fmt.Println(org.Id, org.Name)
	}
}
