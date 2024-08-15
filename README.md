# GoBalance: A Go-based Load Balancer

Welcome to GoBalance, a lightweight and customizable load balancer built in Go. This project demonstrates my proficiency in Go and distributed systems, showcasing my ability to develop robust backend solutions. GoBalance supports multiple load balancing strategies, including Round Robin, Weighted Round Robin, and Sticky Round Robin, and includes health checking functionality for backend servers.

## Features

- **Round Robin**: Distributes incoming requests evenly across all available servers.
- **Weighted Round Robin**: Allows servers to receive traffic proportionally based on their assigned weights.
- **Sticky Round Robin**: Ensures that a client is consistently routed to the same server across requests, enabling session persistence.
- **Health Checks**: Periodically checks the health of backend servers to ensure that only healthy servers receive traffic.

## Getting Started

### Prerequisites

- Go 1.19 or later
- Basic understanding of distributed systems and load balancing concepts

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ParhamMootab/GoBalance.git
   cd GoBalance

2. Build the project:
   
   ```bash
   go build -o gobalance .
   
3. Run the load balancer:
   
   ```bash
   ./gobalance

### Usage

When you run the GoBalance executable, you'll be prompted to configure the load balancer:

Choose a Load Balancing Strategy:

1 for Round Robin
2 for Weighted Round Robin
3 for Sticky Round Robin
Enter Server URLs:

Add the URLs of the backend servers that GoBalance should distribute traffic to.
If you select Weighted Round Robin, you'll also need to specify the weight for each server.
Set Health Check Interval:

Specify the interval (in seconds) at which GoBalance should perform health checks on the backend servers.
Start the Load Balancer:

GoBalance will start listening on port 8080, distributing incoming traffic based on the chosen strategy.

### Example

  ```bash
  $ ./gobalance
  
  ***********************************************************
  *                                                         *
  *      _____       ____        _                          *
  *     / ____|     |  _ \      | |                         *
  *    | |  __  ___ | |_) | __ _| | __ _ _ __   ___ ___     *
  *    | | |_ |/ _ \|  _ < / _`  | |/ _` |  _ \ / __/ _ \    *
  *    | |__| | (_) | |_) | (_| | | (_| | | | | (_|  __/    *
  *     \_____|\___/|____/ \__,_|_|\__,_|_| |_|\___\___|    *
  *                                                         *
  ***********************************************************
  
  Enter the load balancing strategy 
  (1 for Round Robin, 2 for Weighted Round Robin, 3 for Sticky Round Robin): 1
  Enter the server url (Enter 'D' if you've entered all urls): http://localhost:8081
  Enter the server url (Enter 'D' if you've entered all urls): http://localhost:8082
  Enter the server url (Enter 'D' if you've entered all urls): d
  Enter the health check interval in seconds: 10
  
  Load Balancer started at: 8080


### Project Structure


   
   

