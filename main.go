package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	apiConfig := apiConfig{}
	homeHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	// fe
	serveMux.Handle("/app/", apiConfig.middlewareMetricsInc(homeHandler))
	//api
	serveMux.HandleFunc("GET /api/healthz", http.HandlerFunc(handlerHealthz))
	serveMux.HandleFunc("POST /api/validate_chirp", http.HandlerFunc(handlerValidateChirp))
	//admin
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
