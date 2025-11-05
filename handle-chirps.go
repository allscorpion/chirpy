package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/allscorpion/chirpy/internal/database"
	"github.com/google/uuid"
)

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

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();
	type parameters struct {
        Body string `json:"body"`
        UserId uuid.UUID `json:"user_id"`
    }

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

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: getCleanBody(params.Body),
		UserID: params.UserId,
	});

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp");
	}

	type successResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time	`json:"updated_at"`
		Body      string `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	respondWithJSON(w, http.StatusCreated, successResponse{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}