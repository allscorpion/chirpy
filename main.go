package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	ServeHTTP func(w http.ResponseWriter, req *http.Request)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	newVal := cfg.fileserverHits.Add(1);
	cfg.fileserverHits.Store(newVal);
	return next;
}

func (cfg *apiConfig) getMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8");
		w.WriteHeader(200);
		hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
		w.Write([]byte(hits));
	});
}

func (cfg *apiConfig) reset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8");
		w.WriteHeader(200);
		cfg.fileserverHits.Store(0)
	});
}


func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8");
	w.WriteHeader(200);
	w.Write([]byte("OK"));
}


func main() {
	serveMux := http.NewServeMux();
	config := apiConfig{fileserverHits: atomic.Int32{}}

	serveMux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))));
	serveMux.Handle("/healthz", http.HandlerFunc(healthCheck));
	serveMux.Handle("/metrics", config.getMetrics());
	serveMux.Handle("/reset", config.reset());

	server := http.Server{
		Handler: serveMux,
		Addr: ":8080",
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("failed to start server, %v", err)
	}
}