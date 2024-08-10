package LoadBalancer

import (
	"net/http"

	ServerPkg "github.com/ParhamMootab/GoBalance/Server"
	// "log"
	"sync/atomic"
)

type WeightedRoundRobin struct {
	current            uint32
	WeightedServerList []*ServerPkg.Server
}

func (lb *WeightedRoundRobin) getNextServer(r *http.Request, w http.ResponseWriter) *ServerPkg.Server {
	weightSum := len(lb.WeightedServerList)

	nextServer := lb.WeightedServerList[lb.current%uint32(weightSum)]
	atomic.AddUint32(&lb.current, 1)
	if !nextServer.Healthy {
		return lb.getNextServer(r, w)
	}
	return nextServer

}
