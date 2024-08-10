package LoadBalancer

import (
	ServerPkg "github.com/ParhamMootab/GoBalance/Server"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type LoadBalancerStrategy interface {
	getNextServer(r *http.Request, w http.ResponseWriter) *ServerPkg.Server
}

type LoadBalancer struct {
	Strategy            LoadBalancerStrategy
	ServerList          []*ServerPkg.Server
	HealthCheckInterval int
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var server *ServerPkg.Server = lb.Strategy.getNextServer(r, w)
	if server == nil {
		http.Error(w, "No healthy servers available", http.StatusServiceUnavailable)
		return
	}
	log.Printf("Forwarding request to %s", server.URL)
	server.ReverseProxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) HealthCheck() {
	for {
		for _, server := range lb.ServerList {
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

func NewLoadBalancer(serverUrls []string, healthCheckInterval int, weights []int, strategyType int) *LoadBalancer {
	var servers []*ServerPkg.Server

	for index, serverUrl := range serverUrls {
		url, err := url.Parse(serverUrl)
		if err != nil {
			log.Fatalf("Failed to parse server URL: %v", err)
		}
		newServer := &ServerPkg.Server{
			URL:          url,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
			Healthy:      true,
			Weight:       weights[index],
		}
		servers = append(servers, newServer)

	}

	var strategy LoadBalancerStrategy

	switch strategyType {
	case 1:
		strategy = &RoundRobinLoadBalancer{servers: servers}
	case 2:
		var weightedServerList []*ServerPkg.Server
		for index, server := range servers {
			for i := 0; i < weights[index]; i++ {
				weightedServerList = append(weightedServerList, server)
			}
		}
		strategy = &WeightedRoundRobin{WeightedServerList: weightedServerList}
	case 3:
		strategy = &StickyRoundRobin{
			servers:   servers,
			clientMap: make(map[string]*ServerPkg.Server, 0),
		}
	}

	return &LoadBalancer{
		ServerList:          servers,
		HealthCheckInterval: healthCheckInterval,
		Strategy:            strategy,
	}
}
