package main

import (
	"net/http"
	"time"

	"github.com/fergalhk/loadshed-go/netutil"
	"k8s.io/apimachinery/pkg/util/sets"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/do-something-important", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 3)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:    "127.0.0.1:9000",
		Handler: netutil.NewLimitMux(5, mux, sets.New[string]("/livez")),
	}
	panic(srv.ListenAndServe())
}
