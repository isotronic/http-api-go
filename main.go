package main

import "net/http"

func main() {
	var server http.Server
	mux := http.NewServeMux()
	server.Addr = ":8080"
	server.Handler = mux

	server.ListenAndServe()
}