package dulcamara

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type server struct {
	e concurrentEndpoints
	s *http.Server
}

type concurrentEndpoints struct {
	mutex     sync.RWMutex
	endpoints []endpoint
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL
	fmt.Printf("got request with method %s for path %s\n", method, path)
	s.e.mutex.RLock()
	defer s.e.mutex.RUnlock()
	for _, endpoint := range s.e.endpoints {
		if r.URL.String() == endpoint.path && r.Method == endpoint.method {
			fmt.Printf("matched rule %s\n", endpoint.rule)
			io.WriteString(w, endpoint.response)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

var servers = map[string]*server{}

func Deploy(e endpoint) {

	// loop through servers (one could have modified a file with a new port, you would overlook old version if you went straight for map[[prt]])
	// if endpoint found, undeploy()
	// deploy new version
	mockServer, serverExists := servers[e.port]
	if !serverExists {
		mockServer = &server{
			e: concurrentEndpoints{endpoints: []endpoint{e}},
			s: &http.Server{
				Addr: ":" + e.port,
			},
		}
		m := http.NewServeMux()
		m.HandleFunc("/", mockServer.handle)
		mockServer.s.Handler = m
		go func() {
			mockServer.s.ListenAndServe()
		}()
		servers[e.port] = mockServer
	} else {
		mockServer.e.mutex.Lock()
		defer mockServer.e.mutex.Unlock()
		mockServer.e.endpoints = append(mockServer.e.endpoints, e)
	}
	fmt.Printf("added rule %s\n", e.rule)
}

// call on endpoint delete
func Undeploy(ruleName string) {
	for _, server := range servers {
		for i, endpoint := range server.e.endpoints {
			if endpoint.rule == ruleName {
				server.e.mutex.Lock()
				defer server.e.mutex.Unlock()
				server.e.endpoints = remove(server.e.endpoints, i)
				fmt.Printf("Removed rule %s\n", ruleName)
				return
			}
		}
	}
}

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
