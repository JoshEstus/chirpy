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
}

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
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
	}

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	serverMux.Handle("/app/", middlewareLog(apiCfg.middlewareMetricsInc(fileHandler)))

	serverMux.HandleFunc("GET /api/healthz", HealthzHandler)
	serverMux.HandleFunc("POST /api/validate_chirp", ValidateChirpHandler)
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
