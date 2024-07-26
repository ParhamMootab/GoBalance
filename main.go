package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Server struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	servers []*Server
	current uint32
}

func (lb *LoadBalancer) getNextServer() *Server {
	next := atomic.AddUint32(&lb.current, 1)
    return lb.servers[next%uint32(len(lb.servers))]
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.getNextServer()
	log.Printf("Forwarding request to %s", server.URL)
	server.ReverseProxy.ServeHTTP(w, r)
}

func NewLoadBalancer(serverUrls []string) *LoadBalancer {
	var servers []*Server
	for _, serverUrl := range serverUrls {
		url, err := url.Parse(serverUrl)
		if err != nil {
			log.Fatalf("Failed to parse server URL: %v", err)
		}
		servers = append(servers, &Server{
			URL:          url,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		})

	}
	return &LoadBalancer{servers: servers}
}

func main() {
	serverUrls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	lb := NewLoadBalancer(serverUrls)
	http.Handle("/", lb)
	fmt.Println("Load Balancer started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
