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
)

func (p *Proxy) ListenAndServe() {
	l, err := net.Listen("tcp4", fmt.Sprintf("%v:%v", p.Listen.Address, p.Listen.Port))
	if err != nil {
		log.Fatalf("Can't bind to address %v:%v, error: %v", p.Listen.Address, p.Listen.Port, err.Error())
	}
	http.Serve(l, p)
}

func (p *Proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	preProcesStartTime := time.Now()
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
	request.RequestURI = ""
	//iterate randomly over the backend instances until we get a response
	for _, i := range rand.Perm(len(p.ServiceMap[hostHeader].Hosts)) {
		target := p.ServiceMap[hostHeader].Hosts[i]
		//replace the host in the original request's URL
		targetUrl := CompileTargetURL(&target, request.URL)
		request.URL = targetUrl
		preProcessStopTime := time.Now()
		//send the request to the backend service
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			//if we don't get a response, try the next host
			log.Printf("Error connecting to backend host %v:%v - %v", target.Address, target.Port, err.Error())
			continue
		}
		go preProcessHistogram.Observe(float64(preProcessStopTime.Sub(preProcesStartTime).Nanoseconds() / 1000))
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
