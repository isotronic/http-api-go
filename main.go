package main

import (
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
	mux.HandleFunc("GET /api/healthz", apiHealthzHandler)
	mux.HandleFunc("GET /admin/metrics", adminMetricsHandler(&apiCfg))
	mux.Handle("POST /admin/reset", apiCfg.middleWareMetricsReset(http.HandlerFunc(adminResetHandler)))

	server.ListenAndServe()
}

func (cfg *apiConfig) middleWareMetricsInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middleWareMetricsReset(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Swap(0)
		next.ServeHTTP(w, r)
	})
}