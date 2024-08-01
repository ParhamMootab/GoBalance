package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type LoadBalancer struct {
	servers             []*Server
	current             uint32
	HealthCheckInterval int
	WeightedServerList  []*Server
}

// Round Robin
func (lb *LoadBalancer) getNextServer() *Server {

	for range lb.servers {
		nextServer := lb.servers[lb.current%uint32(len(lb.servers))]
		atomic.AddUint32(&lb.current, 1)
		if nextServer.Healthy {
			return nextServer
		}
	}

	log.Fatal("No healthy servers available")
	return nil
}

// Weighted Round Robin
func (lb *LoadBalancer) getNextWeightedServer() *Server {
	weightSum := len(lb.WeightedServerList)

	nextServer := lb.WeightedServerList[lb.current%uint32(weightSum)]
	atomic.AddUint32(&lb.current, 1)
	if !nextServer.Healthy {
		return lb.getNextWeightedServer()
	}
	return nextServer

}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var server *Server
	if len(lb.servers) == len(lb.WeightedServerList) {
		server = lb.getNextServer()
	} else {
		server = lb.getNextWeightedServer()
	}
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
