package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type Server struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
	Healthy      bool
	Weight       int
}

type LoadBalancer struct {
	servers             []*Server
	current             uint32
	HealthCheckInterval int
	WeightedServerList  []*Server
}

// Round Robin
func (lb *LoadBalancer) getNextServer() *Server {

	
	nextServer := lb.servers[lb.current%uint32(len(lb.servers))]
	atomic.AddUint32(&lb.current, 1)
	if nextServer.Healthy {
		
		return nextServer
	}
	
	return lb.getNextServer()
}

// Weighted Round Robin
func (lb *LoadBalancer) getNextWeightedServer() *Server {
	weightSum := len(lb.WeightedServerList)

	next := atomic.AddUint32(&lb.current, 1)
	nextServer := lb.WeightedServerList[next%uint32(weightSum)]
	if !nextServer.Healthy {
		return lb.getNextWeightedServer()
	}
	return nextServer
	

}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	server := lb.getNextServer()
	log.Printf("Forwarding request to %s", server.URL)
	server.ReverseProxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) HealthCheck() {
	for {
		for _, server := range lb.servers {
			resp, err := http.Get(server.URL.String() + "/health")
			if err != nil || resp.StatusCode != http.StatusOK {
				server.Healthy = false
				log.Printf("Server %s is unhealthy", server.URL)
			} else {
				server.Healthy = true
			}

			if resp != nil {
				resp.Body.Close()
			}

		}
		time.Sleep(time.Duration(float64(lb.HealthCheckInterval) * float64(time.Second)))
	}
}

func NewLoadBalancer(serverUrls []string, healthCheckInterval int, weights []int) *LoadBalancer {
	var servers []*Server
	var weightedServerList []*Server
	for index, serverUrl := range serverUrls {
		url, err := url.Parse(serverUrl)
		if err != nil {
			log.Fatalf("Failed to parse server URL: %v", err)
		}
		newServer := &Server{
			URL:          url,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
			Healthy:      true,
			Weight:       weights[index],
		}
		servers = append(servers, newServer)

		for i := 0; i < weights[index]; i++ {
			weightedServerList = append(weightedServerList, newServer)
		}
	}

	return &LoadBalancer{
		servers:             servers,
		HealthCheckInterval: healthCheckInterval,
		WeightedServerList:  weightedServerList,
	}
}

func main() {
	serverUrls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	lb := NewLoadBalancer(serverUrls, 10, []int{1, 1, 1})

	go lb.HealthCheck() // Start health check in a seperate goroutine

	http.Handle("/", lb)
	fmt.Println("Load Balancer started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
