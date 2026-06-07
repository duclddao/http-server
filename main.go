package main

import (
	"database/sql"
	"http-server/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	apiConfig := apiConfig{
		dbQueries: *dbQueries,
		platform:  os.Getenv("PLATFORM"),
	}
	homeHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	// fe
	serveMux.Handle("/app/", apiConfig.middlewareMetricsInc(homeHandler))
	//api
	serveMux.HandleFunc("GET /api/healthz", http.HandlerFunc(handlerHealthz))
	serveMux.HandleFunc("POST /api/users", apiConfig.handlerUser)
	serveMux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirp)
	serveMux.HandleFunc("GET /api/chirps", apiConfig.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerGetChirps)

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
