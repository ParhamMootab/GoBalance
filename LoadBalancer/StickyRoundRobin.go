package LoadBalancer

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	ServerPkg "github.com/ParhamMootab/GoBalance/Server"
)

type StickyRoundRobin struct {
	servers   []*ServerPkg.Server
	clientMap map[string]*ServerPkg.Server
	mu        sync.RWMutex
	current   uint32
}

func (lb *StickyRoundRobin) getNextServer(r *http.Request, w http.ResponseWriter) *ServerPkg.Server {
	clientID, err := r.Cookie("client_id")
	if err != nil {
		clientID = &http.Cookie{
			Name:  "client_id",
			Value: fmt.Sprintf("%d", atomic.AddUint32(&lb.current, 1)),
			Path:  "/",
		}
		http.SetCookie(w, clientID)
	}

	lb.mu.RLock()
	if server, ok := lb.clientMap[clientID.Value]; ok && server.Healthy {
		lb.mu.RUnlock()
		return server
	}
	lb.mu.RUnlock()

	nextServer := lb.servers[lb.current%uint32(len(lb.servers))]
	atomic.AddUint32(&lb.current, 1)

	lb.mu.Lock()
	if !nextServer.Healthy {
		return lb.getNextServer(r, w)
	}
	lb.clientMap[clientID.Value] = nextServer
	lb.mu.Unlock()
	return nextServer
}
