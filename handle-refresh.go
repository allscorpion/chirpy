package main

import (
	"net/http"
	"time"

	"github.com/allscorpion/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();

	refeshToken, err := auth.GetBearerToken(req.Header);

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid token", err);
		return;
	}

	userWithTokenData, err := cfg.dbQueries.GetUserFromRefreshToken(req.Context(), refeshToken);

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to find token record", err);
		return;
	}

	currentTime := time.Now()
	
	if currentTime.After(userWithTokenData.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "token is expired", err);
		return;
	}

	accessToken, err := auth.MakeJWT(userWithTokenData.UserID, cfg.jwt_secret);

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create access token", err);
		return;
	}

	type successResponse struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, successResponse{
		Token: accessToken,
	})
}