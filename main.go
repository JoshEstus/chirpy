package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {

	const filepathRoot = "."
	const port = "8080"

	serverMux := http.NewServeMux()

	apiCfg := apiConfig{}

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	serverMux.Handle("/app/", middlewareLog(apiCfg.middlewareMetricsInc(fileHandler)))

	serverMux.HandleFunc("GET /api/healthz", healthz)
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

func healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
