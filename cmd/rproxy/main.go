package main

import (
	"log"
	"net/http"

	"github.com/mssbg/rproxy/proxy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/live", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
	})
	http.HandleFunc("/ready", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
	})

	go http.ListenAndServe(":2112", nil)

	log.Printf("Starting proxy %v", proxy.P)
	proxy.P.ListenAndServe()
}
