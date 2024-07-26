package main

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/fergalhk/loadshed-go/netutil"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/do-something-important", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 3)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/do-something-intensive", func(w http.ResponseWriter, r *http.Request) {
		for range 100000000 {
			math.Sqrt(rand.Float64())
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:9000",
		Handler: netutil.NewCPULimitMux(1, mux),
		// Handler: netutil.NewLimitMux(5, mux, sets.New[string]("/livez"), time.Second*5),
	}
	log.Print("Starting server")
	panic(srv.ListenAndServe())
}
