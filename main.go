package main

import (
	"log"
	"net/http"
)

func main() {

	const filepathRoot = "."
	const port = "8080"

	serverMux := http.NewServeMux()

	serverMux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())

}
