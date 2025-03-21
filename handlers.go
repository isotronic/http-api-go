package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func apiHealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func apiValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Body string `json:"body"`
	}
	type responseData struct {
		Valid bool `json:"valid"`
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(500)
		return
	}

	resData := responseData{}
	if len([]rune(reqData.Body)) > 140 {
		resData.Error = "Chirp is too long"
	} else {
		resData.Valid = true
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resData.Valid {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	w.Write(data)
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