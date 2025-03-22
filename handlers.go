package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/isotronic/http-go-server/internal/database"
)

func apiHealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func apiChirpsHandler(apiCfg *apiConfig) http.HandlerFunc {return func(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Body string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(500)
		return
	}

	if len([]rune(reqData.Body)) > 140 {
		respondWithError(w, 400, "Your message is too long")
		return
	}

	body := profanityFilter(reqData.Body)
	userID, err := uuid.Parse(reqData.UserID)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		respondWithError(w, 400, "User ID is invalid")
		return
	}

	newChirp := database.CreateChirpParams{UserID: userID, Body: body}
	chirp, err := apiCfg.database.CreateChirp(r.Context(), newChirp)
	if err != nil {
		log.Panicf("Error creating chirp: %v", err)
		respondWithError(w, 500, "Error creating chirp")
		return
	}

	respondWithJSON(w, 201, ChirpResponse(chirp))
}}

func adminResetHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	if apiCfg.platform == "" {
		w.WriteHeader(403)
		return
	}

	err := apiCfg.database.ResetUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %v", err)
		respondWithError(w, 500, "Error deleting users")
		return
	}

	w.WriteHeader(200)
}}

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

func apiCreateUserHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(500)
		return
	}

	if reqData.Email == "" {
		respondWithError(w, 400, "No email was provided")
		return
	}

	user, err := apiCfg.database.CreateUser(r.Context(), reqData.Email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, 500, "Error creating user")
		return
	}

	respondWithJSON(w, 201, UserResponse(user))
}}

func apiGetAllChirpsHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	chirps, err := apiCfg.database.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error fetching chirps: %v", err)
		respondWithError(w, 500, "Error fetching chirps")
		return
	}

	chirpResponse := make([]ChirpResponse, len(chirps))
	for i, chirp := range chirps {
		chirpResponse[i] = ChirpResponse(chirp)
	}

	respondWithJSON(w, 200, chirpResponse)
}}

func apiGetChirpByIdHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	pathParam := r.PathValue("chirpID")
	if pathParam == "" {
		respondWithError(w, 400, "No chirpID provided")
		return
	}

	chirpID, err := uuid.Parse(pathParam)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		respondWithError(w, 500, "ChirpID is invalid")
		return
	}

	chirp, err := apiCfg.database.GetChirpById(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error fetching chirp: %v", err)
		respondWithError(w, 404, "No chirp with that ID exists")
		return
	}
	
	respondWithJSON(w, 200, ChirpResponse(chirp))
}}