package main

import (
	"fmt"
	"log"
)

// TODO verify this
func converge(c client, mail string) {
	last := c.lastKnown()
	if last == *orgs {
		log.Println("desired state == current state. nothing to do")
		return
	}
	if last > *orgs {
		for i := *orgs + 1; i <= last; i++ {
			// TODO get id's of all users and delete them
		}
	} else {
		log.Println("getting list of collectors to use")
		collectorIds := c.publicCollectorIds()

		for o := last; o <= *orgs; o++ {
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
}
