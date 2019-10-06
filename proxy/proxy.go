package proxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/yaml.v2"
)

func (p *Proxy) ListenAndServe() {
	l, err := net.Listen("tcp4", fmt.Sprintf("%v:%v", p.Listen.Address, p.Listen.Port))
	if err != nil {
		log.Fatalf("Can't bind to address %v:%v, error: %v", p.Listen.Address, p.Listen.Port, err.Error())
	}
	//For each Service we will start a Go routine that will provide the index of the next Host to receive a request.
	for _, s := range p.ServiceMap {
		switch p.Strategy {
		case STRATEGY_ROUND_ROBIN:
			c := make(chan int)
			s.NextHost = c
			go func() {
				for {
					for i, _ := range s.Hosts {
						c <- i
					}
				}
			}()

		case STRATEGY_RANDOM:
			c := make(chan int)
			s.NextHost = c
			go func() {
				length := len(s.Hosts)
				for {
					c <- rand.Intn(length)
				}
			}()
		default:
			log.Fatalf("Unknown strategy %v", p.Strategy)
		}
	}
	http.Serve(l, p)
}

//When a request arrives, this method is started in a new Go routine.
func (p *Proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	preProcessStartTime := time.Now()
	//This property holds the value of the "Host" header
	hostHeader := request.Host
	if hostHeader == "" {
		log.Printf("Can't get Host header %v", request.Host)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	//check if we have that service
	if p.ServiceMap[hostHeader] == nil {
		log.Printf("Can't find service for host %v", hostHeader)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	//For faster processing, we will avoid copying the original request and we will just re-use it to send it
	//to the backend Host.
	//To do that, however, we have to zero out this field, since it is not allowed to be set
	//when sending it to http.Client
	request.RequestURI = ""

	//When a request arrives, we will try to send it to a Host.
	//If we fail to send it and/or get a response, we will try the "next" one. What's "next" depends on the strategy.
	//For any given request we don't want to try any Host more than once.
	testedHosts := map[int]bool{} //this the closest to a Set in Go
	var service = p.ServiceMap[hostHeader]
	for len(testedHosts) < len(service.Hosts) { //when the sizes are equal, we have tried all Hosts
		//obtain the index of the Host to which we are going to send the request
		i := <-service.NextHost
		if testedHosts[i] {
			//we've already tried this host, so try another one
			continue
		}
		testedHosts[i] = true
		target := service.Hosts[i]
		//replace the host in the original request's URL
		request.URL = CompileTargetURL(&target, request.URL)
		preProcessStopTime := time.Now()
		//send the request to the backend service
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			//if we don't get a response, try the next host
			log.Printf("Error connecting to backend host %v:%v - %v", target.Address, target.Port, err.Error())
			continue
		}
		go preProcessHistogram.Observe(float64(preProcessStopTime.Sub(preProcessStartTime).Nanoseconds() / 1000))
		postProcessStartTime := time.Now()
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			log.Printf("can't read response body")
			writer.WriteHeader(http.StatusBadGateway)
			return
		}
		//send back the response's headers
		for k, h := range response.Header {
			for _, v := range h {
				writer.Header().Add(k, v)
			}
		}
		//send back the response's status code
		writer.WriteHeader(response.StatusCode)
		//send the response's body
		writer.Write(body)
		postProcessStopTime := time.Now()
		go postProcessHistogram.Observe(float64(postProcessStopTime.Sub(postProcessStartTime).Nanoseconds() / 1000))
		return
	}
	log.Printf("could not find live host")
	//if we exit the loop, it means none of the backend hosts are reachable
	writer.WriteHeader(http.StatusBadGateway)
	return
}

//Replaces the host address and port in URL instances with those from a Host instance
func CompileTargetURL(host *Host, originalUrl *url.URL) *url.URL {
	scheme := originalUrl.Scheme
	if scheme == "" {
		scheme = "http"
	}
	return &url.URL{
		Scheme:     scheme,
		Opaque:     originalUrl.Opaque,
		User:       originalUrl.User,
		Host:       fmt.Sprintf("%v:%v", host.Address, host.Port),
		Path:       originalUrl.Path,
		RawPath:    originalUrl.RawPath,
		ForceQuery: originalUrl.ForceQuery,
		RawQuery:   originalUrl.RawQuery,
		Fragment:   originalUrl.Fragment,
	}
}

//Used when "printing" the Proxy struct
func (p *Proxy) String() string {
	yamlBytes, err := yaml.Marshal(p)
	if err != nil {
		return "Error serializing Proxy struct"
	}
	return string(yamlBytes)
}
