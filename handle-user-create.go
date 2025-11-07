package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/allscorpion/chirpy/internal/auth"
	"github.com/allscorpion/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();
	type parameters struct {
        Email string `json:"email"`
		Password string `json:"password"`
    }

	decoder := json.NewDecoder(req.Body);
	var params parameters;
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 500, "Something went wrong", err);
		return;
	}

	hash, err := auth.HashPassword(params.Password);

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse password", err);
		return;
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	});

	if err != nil {
		respondWithError(w, 500, "failed to create user", err);
		return;
	}

	type successResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string `json:"email"`
	}

	respondWithJSON(w, 201, successResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	});
}