package dulcamara

import (
	"fmt"
	"io"
	"net/http"
)

type server struct {
	endpoints []endpoint
	s         *http.Server
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL
	fmt.Printf("got request with method %s for path %s\n", method, path)
	//TODO concurrency
	for _, endpoint := range s.endpoints {
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
			endpoints: []endpoint{e},
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
		//TODO concurrency
		mockServer.endpoints = append(mockServer.endpoints, e)
	}
	fmt.Printf("added rule %s\n", e.rule)
}

// call on endpoint delete
func undeploy(e endpoint) {}
