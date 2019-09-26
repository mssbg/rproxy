package echo_server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
)

type Echo struct {
	Address      string
	Counter      int
	CounterMutex sync.Mutex
}

type EchoResponse struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body"`
}

func (e *Echo) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	e.CounterMutex.Lock()
	e.Counter++
	e.CounterMutex.Unlock()
	go log.Printf("Served requests count %v ", e.Counter)
	body, _ := ioutil.ReadAll(request.Body)
	r := EchoResponse{
		Method: request.Method,
		URL:    request.URL.String(),
		Body:   string(body),
	}
	j, _ := json.MarshalIndent(r, "", "  ")
	writer.Write(j)
}

func (e *Echo) ListenAndServe() {
	l, err := net.Listen("tcp4", e.Address)
	if err != nil {
		log.Fatalf("Can't bind to address %v, error: %v", e.Address, err.Error())
	}
	http.Serve(l, e)
}
