package main

import (
	"fmt"
	"net/http"
)

func apiHealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func adminResetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func adminMetricsHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)

	body := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, 
		apiCfg.fileServerHits.Load())
	
		w.Write([]byte(body))
}}