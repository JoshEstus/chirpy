package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/joshestus/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	serverMux := http.NewServeMux()

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	serverMux.Handle("/app/", middlewareLog(apiCfg.middlewareMetricsInc(fileHandler)))

	// API
	serverMux.HandleFunc("GET /api/healthz", HealthzHandler)

	// User
	serverMux.HandleFunc("POST /api/users", apiCfg.CreateUserHandler)

	// Chirp
	serverMux.HandleFunc("POST /api/chirps", apiCfg.AddChirpHandler)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.GetAllChirpsHandler)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpHandler)

	// Admin
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.FileServerHitsResetHandler)

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())

}
