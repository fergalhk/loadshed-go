package netutil

import (
	"context"
	"log"
	"net/http"
)

type LimitMux struct {
	semCh   chan struct{}
	handler http.Handler
}

func NewLimitMux(limit uint, handler http.Handler) http.Handler {
	return &LimitMux{
		semCh:   make(chan struct{}, limit),
		handler: handler,
	}
}

func (l *LimitMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := l.acquire(r.Context()); err != nil {
		log.Printf("Error acquiring semaphore: %s", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer l.release()
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
