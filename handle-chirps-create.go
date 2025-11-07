package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/allscorpion/chirpy/internal/auth"
	"github.com/allscorpion/chirpy/internal/database"
	"github.com/google/uuid"
)

type ChirpParsed struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time	`json:"updated_at"`
		Body      string `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
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

func (cfg *apiConfig) handleChirpsCreate(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();
	type parameters struct {
        Body string `json:"body"`
    }

	decoder := json.NewDecoder(req.Body);
	var params parameters;
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 500, "Something went wrong", err);
		return;
	}

	bearerToken, err := auth.GetBearerToken(req.Header);

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", err);
		return;
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwt_secret);

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", err);
		return;
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil);
		return;
	}

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: getCleanBody(params.Body),
		UserID: userId,
	});

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp", err);
		return;
	}

	respondWithJSON(w, http.StatusCreated, ChirpParsed{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}