package main

import "net/http"

func main() {
	var server http.Server
	mux := http.NewServeMux()
	server.Addr = ":8080"
	server.Handler = mux

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	server.ListenAndServe()
}