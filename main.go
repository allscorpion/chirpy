package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/allscorpion/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK);
	content := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load());
	w.Write([]byte(content))
}


func (cfg *apiConfig) reset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200);
	w.Write([]byte("Hits reset to 0"))
}


func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8");
	w.WriteHeader(200);
	w.Write([]byte("OK"));
}



func respondWithError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.WriteHeader(statusCode);
	type errorResponse struct {
		Error string `json:"error"`
	}

	errorResp := errorResponse{
		Error: errMsg,
	}

	data, err := json.Marshal(errorResp);

	if err != nil {
		w.Write([]byte(fmt.Sprintf("unable to marshal json: %v", data)))
		return;
	}

	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode);
	data, err := json.Marshal(payload);

	if err != nil {
		w.Write([]byte(fmt.Sprintf("unable to marshal json: %v", data)))
		return;
	}

	w.Write(data)
}


func generateSuccessBody(w http.ResponseWriter, cleanedBody string) {
	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	resp := successResponse{
		CleanedBody: cleanedBody,
	}

	respondWithJSON(w, 200, resp);
}

func getCleanBody(body string) string {
	invalidWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	};
	words := strings.Fields(body);

	for i, word := range words {
		parsedWord := strings.ToLower(word);

		if _, exists := invalidWords[parsedWord]; exists {
			words[i] = "****";
		}
	}

	return strings.Join(words, " ");
}

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
        Body string `json:"body"`
    }

	w.Header().Set("Content-Type", "application/json");
	decoder := json.NewDecoder(req.Body);
	var params parameters;
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 500, "Something went wrong");
		return;
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long");
		return;
	}

	generateSuccessBody(w, getCleanBody(params.Body));
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

	dbQueries := database.New(db)

	serveMux := http.NewServeMux();
	config := apiConfig{
		fileserverHits: atomic.Int32{}, 
		dbQueries: dbQueries,
	};

	serveMux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))));
	serveMux.HandleFunc("GET /api/healthz", healthCheck);
	serveMux.HandleFunc("GET /admin/metrics", config.handlerMetrics);
	serveMux.HandleFunc("POST /admin/reset", config.reset);
	serveMux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	server := http.Server{
		Handler: serveMux,
		Addr: ":8080",
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("failed to start server, %v", err)
	}
}

