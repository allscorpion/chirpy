package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/allscorpion/chirpy/internal/auth"
	"github.com/allscorpion/chirpy/internal/database"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();

	refeshToken, err := auth.GetBearerToken(req.Header);

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid token", err);
		return;
	}

	err = cfg.dbQueries.SetRevokedAt(req.Context(), database.SetRevokedAtParams{
		Token: refeshToken,
		RevokedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
	});

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update record", err);
		return;
	} 

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}