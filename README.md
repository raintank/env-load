### env-load

manage artificial, but realistic loads on a raintank litmus environment.
so that you can proactively test and validate the entire stack for a given workload.


### operation

```
env-load -host=https://grafana-endpoint -auth=admin:admin -orgs=4 load
env-load -host=https://grafana-endpoint -auth=admin:admin status
env-load -host=https://grafana-endpoint -auth=admin:admin clean
```

requires authentication as grafana super admin.

`load` mass-creates users with 1 org each, 4 endpoints per org and 3 monitors per endpoint (ping, dns, http)
`status` gives an overview of existing orgs, users and endpoints
`clean` removes all fake orgs, users, endpoints and monitors again.

