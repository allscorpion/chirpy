package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/allscorpion/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfigEnv struct {
	platform string
	jwt_secret string
	polka_key string
}

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	env apiConfigEnv
}

type customUser struct {
	ID        	uuid.UUID `json:"id"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Email     	string `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func getEnvVariable(key string) string {
	val := os.Getenv(key);

	if val == "" {
		log.Fatalf("%v is not defined", key);
	}

	return val;
}

func main() {
	err := godotenv.Load();

	if err != nil {
		log.Fatal("failed to load env");
	}

	dbURL := getEnvVariable("DB_URL");

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("failed to open sql connection");
	}

	dbQueries := database.New(db);

	platform := getEnvVariable("PLATFORM");
	jwt_secret := getEnvVariable("JWT_SECRET");
	polka_key := getEnvVariable("POLKA_KEY");

	serveMux := http.NewServeMux();
	
	config := apiConfig{
		fileserverHits: atomic.Int32{}, 
		dbQueries: dbQueries,
		env: apiConfigEnv{
			platform: platform,
			jwt_secret: jwt_secret,
			polka_key: polka_key,
		},
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
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", config.handleChirpsDelete)
	serveMux.HandleFunc("POST /api/refresh", config.handleRefresh)
	serveMux.HandleFunc("POST /api/revoke", config.handleRevoke)
	serveMux.HandleFunc("POST /api/polka/webhooks", config.handlePolkaWebook)

	server := http.Server{
		Handler: serveMux,
		Addr: ":8080",
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("failed to start server, %v", err)
	}
}

