package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/allscorpion/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
	jwt_secret string
}

func main() {
	err := godotenv.Load();

	if err != nil {
		log.Fatal("failed to load env");
	}

	dbURL := os.Getenv("DB_URL");

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("failed to open sql connection");
	}

	dbQueries := database.New(db);

	platform := os.Getenv("PLATFORM");
	jwt_secret := os.Getenv("PLATFORM");

	serveMux := http.NewServeMux();
	config := apiConfig{
		fileserverHits: atomic.Int32{}, 
		dbQueries: dbQueries,
		platform: platform,
		jwt_secret: jwt_secret,
	};

	serveMux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))));
	serveMux.HandleFunc("GET /api/healthz", healthCheck);
	serveMux.HandleFunc("GET /admin/metrics", config.handlerMetrics);
	serveMux.HandleFunc("POST /admin/reset", config.reset);
	serveMux.HandleFunc("POST /api/users", config.handleCreateUser)
	serveMux.HandleFunc("PUT /api/users", config.handleUserUpdate)
	serveMux.HandleFunc("POST /api/login", config.handleLogin)
	serveMux.HandleFunc("POST /api/chirps", config.handleChirpsCreate)
	serveMux.HandleFunc("GET /api/chirps", config.handleChirpsGetAll)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", config.handleChirpsGet)
	serveMux.HandleFunc("POST /api/refresh", config.handleRefresh)
	serveMux.HandleFunc("POST /api/revoke", config.handleRevoke)

	server := http.Server{
		Handler: serveMux,
		Addr: ":8080",
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("failed to start server, %v", err)
	}
}

