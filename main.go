package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/isotronic/http-go-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	database *database.Queries
}

func main() {
	godotenv.Load()
	var server http.Server
	var apiCfg apiConfig
	mux := http.NewServeMux()
	server.Addr = ":8080"
	server.Handler = mux

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	apiCfg.database = database.New(db)

	mux.Handle("/app/", apiCfg.middleWareMetricsInt(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", apiHealthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiValidateChirpHandler)
	mux.HandleFunc("POST /api/users", apiCreateUserHandler(&apiCfg))

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