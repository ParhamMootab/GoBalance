package main

import (
	"fmt"
	"log"
	"net/http"
	
)

// Mitigate the Favicon Issue -> In  round robin each request is sent to the next server in order. 
//		If you want to send the favicon request of the same user to the same server as the main request
//		do it in the sticky round robin Algo not the round robin.
		
// Add Weighted round robin implementation [X]
// Add sticky round robin []
// Check unhealthy servers edge case
// Create the cli API

func main() {
	serverUrls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	lb := NewLoadBalancer(serverUrls, 10, []int{1, 3, 1})

	go lb.HealthCheck() // Start health check in a seperate goroutine

	http.Handle("/", lb)
	fmt.Println("Load Balancer started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
