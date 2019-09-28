# rproxy

A simple reverse proxy implementation. 

It can serve multiple "services", each backed by one or more hosts/instances.
The correct service is identified by the "Host" header in the request.

Requests are relayed to a random backend host.
If that host doesn't respond, another will be tried.

The order in which the backend hosts of particular service will be contacted is
randomized for each request.

 
## Building
To build this project, use Go Language v.1.12 or later

The project has two executable binaries - "rproxy", the reverse proxy itself and 
"echo" a test web service that echoes back any request sent to it. 

To build a locally executable rproxy:

`GO111MODULE=on go build -o rproxy ./cmd/rproxy/`

from the root of the project.

To build the echo server:

`GO111MODULE=on go build -o echo ./cmd/echo/`

### Docker image
Included are a Dockerfile and a build script "build.sh"

To build the docker image, just run the build.sh script.
  
### Kubernetes Helm chart
There is a helm chart provided in the ./helm sub directory.

To deploy it, run:

`cd ./helm && helm install rproxy`

## Running 

To run the "rproxy" simply run the executable:

`./rproxy`

it expects to find the configuration file "config.yml" in the current path.

To run the "echo" server:

`./echo ip_address:port`

## Metrics
The service exposes runtime metrics using the "prometheus" library. The metrics
are accessible via HTTP under "/metrics" on port 2112. 

Beside the standard built-in metrics, the service exposes a couple of Prometheus
Histograms measuring the delay introduced by the reverse proxy. The "pre_process"
histogram measures the delay between receiving a request and successfully relaying
it to a backend host. This will include any potential timeout delays. The 
"post_process" histogram measures the delay between receiving a response from a 
backend host and relaying it back to the client.   
 
 These metrics can be used as SLI. For example the SLI can be defined as:
 
 `99% of all requests should have pre- and post- process times under 100 microseconds`
 
 or whatever threshold is deemed "normal" in a particular setup.
 
  