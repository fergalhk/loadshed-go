package netutil

import (
	"context"
	"log"
	"net/http"

	"k8s.io/apimachinery/pkg/util/sets"
)

type LimitMux struct {
	semCh       chan struct{}
	handler     http.Handler
	bypassPaths sets.Set[string]
}

func NewLimitMux(limit uint, handler http.Handler, bypassPaths sets.Set[string]) http.Handler {
	return &LimitMux{
		semCh:       make(chan struct{}, limit),
		handler:     handler,
		bypassPaths: bypassPaths,
	}
}

func (l *LimitMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if l.shouldBypass(r) {
		l.forwardRequest(w, r)
		return
	}

	if err := l.acquire(r.Context()); err != nil {
		log.Printf("Error acquiring semaphore: %s", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	defer l.release()
	l.forwardRequest(w, r)
}

func (l *LimitMux) shouldBypass(r *http.Request) bool {
	return l.bypassPaths.Has(r.URL.Path)
}

func (l *LimitMux) forwardRequest(w http.ResponseWriter, r *http.Request) {
	l.handler.ServeHTTP(w, r)
}

func (l *LimitMux) acquire(ctx context.Context) error {
	select {
	case l.semCh <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *LimitMux) release() {
	<-l.semCh
}
