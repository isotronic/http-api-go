package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/isotronic/http-go-server/internal/auth"
	"github.com/isotronic/http-go-server/internal/database"
)

func apiHealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func apiPostChirpsHandler(apiCfg *apiConfig) http.HandlerFunc {return func(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(token, apiCfg.tokenSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err = decoder.Decode(&reqData)
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
	newChirp := database.CreateChirpParams{UserID: userID, Body: body}
	chirp, err := apiCfg.database.CreateChirp(r.Context(), newChirp)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
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

	err = apiCfg.database.ResetChirps(r.Context())
	if err != nil {
		log.Printf("Error deleting chirps: %v", err)
		respondWithError(w, 500, "Error deleting chirps")
		return
	}

	err = apiCfg.database.ResetRefreshTokens(r.Context())
	if err != nil {
		log.Printf("Error deleting refresh tokens: %v", err)
		respondWithError(w, 500, "Error deleting refresh tokens")
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
		Password string `json:"password"`
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
	if reqData.Password == "" {
		respondWithError(w, 400, "No password was provided")
		return
	}

	passHash, err := auth.HashPassword(reqData.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, 500, "Error hashing password")
	}

	createUserParams := database.CreateUserParams{Email: reqData.Email, HashedPassword: passHash}
	user, err := apiCfg.database.CreateUser(r.Context(), createUserParams)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, 500, "Error creating user")
		return
	}

	newUser := UserResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	respondWithJSON(w, 201, newUser)
}}

func apiUpdateUserHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(token, apiCfg.tokenSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err = decoder.Decode(&reqData)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(500)
		return
	}

	if reqData.Email == "" {
		respondWithError(w, 400, "No email was provided")
		return
	}
	if reqData.Password == "" {
		respondWithError(w, 400, "No password was provided")
		return
	}

	passHash, err := auth.HashPassword(reqData.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, 500, "Error hashing password")
	}

	updateUserParams := database.UpdateUserParams{ID: userID, Email: reqData.Email, HashedPassword: passHash}
	user, err := apiCfg.database.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		respondWithError(w, 500, "Error updating user")
		return
	}

	updatedUser := UserResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	respondWithJSON(w, 200, updatedUser)
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

func apiLoginHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(500)
		return
	}

	if reqData.Email == "" || reqData.Password == "" {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	user, err := apiCfg.database.GetUserByEmail(r.Context(),reqData.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(reqData.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	expiresIn := time.Duration(60 * 60) * time.Second
	jwt, err := auth.MakeJWT(user.ID, apiCfg.tokenSecret, expiresIn)
	if err != nil {
		log.Printf("Error making token: %v", err)
		w.WriteHeader(500)
	}

	refresh, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error making refresh token: %v", err)
		w.WriteHeader(500)
	}

	insertParams := database.InsertRefreshTokenParams{
		Token: refresh,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}
	_, err = apiCfg.database.InsertRefreshToken(r.Context(), insertParams)
	if err != nil {
		log.Printf("Error inserting refresh token: %v", err)
		w.WriteHeader(500)
		return
	}

	response := LoginResponse{
		AccessToken: jwt,
		RefreshToken: refresh,
	}
	respondWithJSON(w, 200, response)
}}

func apiRefreshHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	refreshEntry, err := apiCfg.database.GetUserFromRefreshToken(r.Context(), refresh)
	if err != nil || refreshEntry.ExpiresAt.Before(time.Now()) || refreshEntry.RevokedAt.Valid  {
		respondWithError(w, 401, "Invalid or expired token")
		return
	}

	jwt, err := auth.MakeJWT(refreshEntry.UserID, apiCfg.tokenSecret, time.Duration(60 * 60) * time.Second)
	if err != nil {
		log.Printf("Error making JWT: %v", err)
		w.WriteHeader(500)
		return
	}

	response := RefreshResponse{
		AccessToken: jwt,
	}
	respondWithJSON(w, 200, response)
}}

func apiRevokeHandler(apiCfg *apiConfig) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
	}

	_, err = apiCfg.database.RevokeRefreshToken(r.Context(), refresh)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		w.WriteHeader(500)
	}
	w.WriteHeader(204)
}}