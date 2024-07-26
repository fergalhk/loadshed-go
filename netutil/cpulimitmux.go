package netutil

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fergalhk/loadshed-go/netutil/stat"
	"github.com/samber/lo"
)

type CPULimitMux struct {
	mu                *sync.RWMutex
	maxLoad, currLoad float64
	handler           http.Handler
	procStats         *stat.ProcessStats
}

func NewCPULimitMux(max float64, handler http.Handler) http.Handler {
	cl := &CPULimitMux{
		mu:        new(sync.RWMutex),
		maxLoad:   max,
		handler:   handler,
		procStats: lo.Must(stat.NewProcessStats()),
	}
	go cl.runStatsRefresher()
	return cl
}

func (c *CPULimitMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !c.isWithinLimits() {
		log.Printf("Rejecting request, CPU of %f is over max of %f", c.currLoad, c.maxLoad)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	c.handler.ServeHTTP(w, r)
}

func (c *CPULimitMux) isWithinLimits() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currLoad <= c.maxLoad
}

func (c *CPULimitMux) runStatsRefresher() {
	tick := time.NewTicker(time.Millisecond * 200)
	defer tick.Stop()
	for range tick.C {
		c.refreshStats()
	}
}

func (c *CPULimitMux) refreshStats() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.currLoad = lo.Must(c.procStats.CPULoadSinceLastCall())
}
