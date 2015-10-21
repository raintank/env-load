package main

import (
	"flag"
	"log"

	"github.com/grafana/grafana-api-golang-client"
)

var orgs = flag.Int("orgs", 100, "number of orgs to create")
var email = flag.String("email", "", "email addresses for alerting emails. empty means notifications disabled")
var host = flag.String("host", "http://localhost/", "https://which.raintank.instance/")
var auth = flag.String("auth", "admin:admin", "authentication string. either 'user:pass' or 'apikey'")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("must have either 'load' or 'clean' operation argument")
	}

	if *host == "" {
		log.Fatal("need a host to connect to")
	}
	if *auth == "" {
		log.Fatal("need an authentication string")
	}
	client, err := gapi.New(*auth, *host)
	if err != nil {
		log.Fatal(err)
	}

	op := args[0]
	switch op {
	case "load":
		if *orgs <= 0 {
			log.Fatal("number of orgs must >= 1")
		}
		load(client, *email)
	case "clean":
		clean(client)
	case "status":
		status(client)
	default:
		log.Fatalf("no such operation %q", op)
	}
}
