package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/grafana/pkg/api/dtos"
	m "github.com/grafana/grafana/pkg/models"
)

var orgs = flag.Int("orgs", 100, "number of orgs to create")
var email = flag.String("email", "", "who will get alerting emails")
var host = flag.String("host", "http://localhost/", "https://which.raintank.instance/")
var auth = flag.String("auth", "", "authentication string. either 'user:pass' or 'apikey'")

func main() {
	flag.Parse()

	if *orgs <= 0 {
		log.Fatal("number of orgs must >= 1")
	}
	if *host == "" {
		log.Fatal("need a host to connect to")
	}
	if *auth == "" {
		log.Fatal("need an authentication string")
	}

	client, err := New(*auth, *host)
	if err != nil {
		log.Fatal(err)
	}

	existingOrgs, err := client.Orgs()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(existingOrgs)

	collectors, err := client.Collectors()
	if err != nil {
		log.Fatal(err)
	}
	collectorIds := make([]int64, 0)
	for _, coll := range collectors {
		collectorIds = append(collectorIds, int64(coll.Id))
	}
	fmt.Println(collectorIds)

	for o := 1; o <= *orgs; o++ {
		user := fmt.Sprintf("fake_user_%d", o)
		pass := fmt.Sprintf("fake_pass_%d", o)
		mail := fmt.Sprintf("fake_user_%d@example.org", o)
		org := fmt.Sprintf("fake_org_%d", o)
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
		subClient, err := New(fmt.Sprintf("%s:%s", user, pass), *host)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(">> creating org", org)
		if err = subClient.NewOrg(org); err != nil {
			log.Fatal(err)
		}
		for e := 1; e <= 4; e++ {
			settings := m.AddEndpointCommand{
				OrgId: 10,
				Name:  fmt.Sprintf("fake_org_%d_endpoint_%d", o, e),
				Tags:  make([]string, 0),
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
