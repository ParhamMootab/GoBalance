package main

import (
	"fmt"
	"net/url"
	"github.com/ParhamMootab/GoBalance/LoadBalancer"
	"log"
	"net/http"
)

// Add Weighted round robin implementation [X]
// Add sticky round robin [X]
// Check unhealthy servers edge case{
// + Sticky round robin doesnt work well when one or multiple servers are down}
// Create the cli API

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}

	// Check if the URL has a valid scheme and host
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}

func main() {

	fmt.Println(`
	***********************************************************
	*                                                         *
	*      _____       ____        _                          *
	*     / ____|     |  _ \      | |                         *
	*    | |  __  ___ | |_) | __ _| | __ _ _ __   ___ ___     *
	*    | | |_ |/ _ \|  _ < / _` + "  | |/ _`" + ` |  _ \ / __/ _ \    *
	*    | |__| | (_) | |_) | (_| | | (_| | | | | (_|  __/    *
	*     \_____|\___/|____/ \__,_|_|\__,_|_| |_|\___\___|    *
	*                                                         *
	***********************************************************`)

	var strategyType int
	for strategyType > 3 || strategyType < 1 {
		fmt.Println("Enter the load balancing strategy \n(1 for Round Robin, 2 for Weighted Round Robin, 3 for Sticky Round Robin): ")
		fmt.Scanln(&strategyType)
	}

	var serverUrls []string
	var weights []int
	doneFlag := false
	for !doneFlag {
		var url string
		fmt.Println("Enter the server url (Enter 'D' if you've entered all urls): ")
		fmt.Scanln(&url)
		if url == "d" {
			doneFlag = true
		}
		for !isValidURL(url) {
			fmt.Println("Wrong format. Try again: ")
			fmt.Scanln(&url)
		}
		serverUrls = append(serverUrls, url)
		var weight int
		if strategyType == 2 {
			for weight < 1 {
				fmt.Println("Enter the weight for this server: ")
				fmt.Scanln(&weight)
			}
			weights = append(weights, weight)
		}
		
	}
	var healthCheckInterval int
	for healthCheckInterval <= 0 {
		fmt.Println("Enter the health check interval in seconds: ")
		fmt.Scanln(&healthCheckInterval)
	}



	// serverUrls := []string{
	// 	"http://localhost:8081",
	// 	"http://localhost:8082",
	// 	"http://localhost:8083",
	// }

	lb := LoadBalancer.NewLoadBalancer(serverUrls, healthCheckInterval, weights, strategyType)

	go lb.HealthCheck() // Start health check in a seperate goroutine

	http.Handle("/", lb)
	fmt.Println("Load Balancer started at: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
