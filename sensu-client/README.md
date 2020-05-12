Small tool to simulate running a sensu client in standalone mode with AMQP1 addresses
# Build

`go build -v`

# CLI

```
./sensu-client [OPTION] 
  -address string
    	amqp1 url and address (default "amqp://0.0.0.0:5672")
  -command string
    	command to run
  -file string
    	filepath to script
  -interval int
    	interval at which to execute script in seconds (default 15)
  -target string
    	target address e.g. /example-queue-name
```

# Example
This example creates a collectd metric and sends it to the AMQP1 address with the intent of it being stored in Prometheus.
Create a file `metric.py` and paste tho following contents into it:

```
import json, time, random
x = [
    {
            "values": [random.random() * 1000],
            "dstypes": ["gauge"],
            "dsnames": ["value"],
            "time": time.time(),
            "interval": 10.000,
            "host": "localhost",
            "plugin": "cpu",
            "plugin_instance": "7",
            "type": "percent",
            "type_instance": "idle"
    }
]


y = json.dumps(x)
print(y)
```
Use the sensu client to run this script and push to the AMQP:
`./sensu-client -command "python metric.py" -target /collectd/telemetry -interval 1`

The above command posts to the address `amqp://0.0.0.0:5672/collectd/telemetry`. Thus, there 
must be a qpid dispatch router at address `0.0.0.0:5672` and some other application listening
at address `/collectd/temeletry`, such as Red Hat's smart gateway.