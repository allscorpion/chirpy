package main

import (
	"encoding/json"
	"net/http"

	"github.com/allscorpion/chirpy/internal/auth"
	"github.com/allscorpion/chirpy/internal/database"
)

func (cfg *apiConfig) handleUserUpdate(w http.ResponseWriter, req *http.Request) {
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

	hashedPassword, err := auth.HashPassword(params.Password);

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password", err);
		return;
	}

	user, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		ID: userId,
		Email: params.Email,
		HashedPassword: hashedPassword,
	});

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user", err);
		return;
	}

	respondWithJSON(w, http.StatusOK, customUser{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}