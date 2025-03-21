package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	var server http.Server
	var apiCfg apiConfig
	mux := http.NewServeMux()
	server.Addr = ":8080"
	server.Handler = mux

	mux.Handle("/app/", apiCfg.middleWareMetricsInt(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		body := fmt.Sprintf("Hits: %v", apiCfg.fileServerHits.Load())
		w.Write([]byte(body))
	})

	server.ListenAndServe()
}

func (cfg *apiConfig) middleWareMetricsInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}