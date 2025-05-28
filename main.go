package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {

	const filepathRoot = "."
	const port = "8080"

	serverMux := http.NewServeMux()

	apiCfg := apiConfig{}

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	serverMux.Handle("/app/", middlewareLog(apiCfg.middlewareMetricsInc(fileHandler)))

	serverMux.HandleFunc("GET /api/healthz", healthz)
	serverMux.HandleFunc("GET /api/metrics", apiCfg.FileServerHitsHandler)
	serverMux.HandleFunc("POST /api/reset", apiCfg.FileServerHitsResetHandler)

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())

}

func healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) FileServerHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf("Hits: %v", cfg.fileServerHits.Load())
	w.Write([]byte(res))
}

func (cfg *apiConfig) FileServerHitsResetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
