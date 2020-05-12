This is a small tool that pulls the json of a Grafana dashboard from a Grafana API server, formats the data to be importable into other Grafana instances without error and stores the data into a specified location. This eases the process of saving and distributing dashboards.
# Build

`go build -v`

# Run

```
usage: ./grafana-save -uid <dashboard uid> [options]
  -address string
    	Grafana API server address (default "localhost:3000")
  -output string
    	Location for dashboard output (default "dashboard.yml")
  -uid string
    	dashboard UID to pull
```

Example:
`./grafana-save -uid eSUyoKjWz -address localhost:3001 -output /tmp/my-dashboard.yml`


# Future Work

Add option to format the data as the Grafana operator dashboard type for kubernetes and OpenShift

