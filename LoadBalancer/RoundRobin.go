package LoadBalancer

import (
	"log"
	"net/http"
	"sync/atomic"

	ServerPkg "github.com/ParhamMootab/GoBalance/Server"
)

type RoundRobinLoadBalancer struct {
	servers []*ServerPkg.Server
	current uint32
}

// Round Robin
func (lb *RoundRobinLoadBalancer) getNextServer(_ *http.Request, _ http.ResponseWriter) *ServerPkg.Server {

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
