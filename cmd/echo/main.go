package main

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/mssbg/rproxy/echo_server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("USAGE: ./echo_server host:port")
	}
	address := os.Args[1]
	parts := strings.Split(address, ":")
	if len(parts) < 2 {
		log.Fatalf("Can't split %v", address)
	}
	e := echo_server.Echo{
		Address:      address,
		Counter:      0,
		CounterMutex: sync.Mutex{},
	}
	e.ListenAndServe()
}
