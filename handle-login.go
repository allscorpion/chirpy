package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/allscorpion/chirpy/internal/auth"
	"github.com/allscorpion/chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
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

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email);

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err);
		return;
	}

	same, err := auth.CheckPasswordHash(params.Password, user.HashedPassword);

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err);
		return;
	}

	if !same {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err);
		return;
	}

	jwt_token, err := auth.MakeJWT(user.ID, cfg.env.jwt_secret);

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create token", err);
		return;
	}

	refreshToken, err := auth.MakeRefreshToken();

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh token", err);
		return;
	}

	_, err = cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	});

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh token record", err);
		return;
	}

	type successResponse struct {
		customUser
		Token     string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, successResponse{
		Token: jwt_token,
		RefreshToken: refreshToken,
		customUser: customUser{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	});
}