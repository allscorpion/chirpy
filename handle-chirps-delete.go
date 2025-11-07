package main

import (
	"net/http"

	"github.com/allscorpion/chirpy/internal/auth"
	"github.com/allscorpion/chirpy/internal/database"
	"github.com/google/uuid"
)



func (cfg *apiConfig) handleChirpsDelete(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();
	
	chirpId, err := uuid.Parse(req.PathValue("chirpID"));

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirpID", err);
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

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpId);

	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to get chirp", err);
		return;
	}

	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "not authorized", err);
		return;
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), database.DeleteChirpParams{
		UserID: userId,
		ID: chirpId,
	});

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete chirp", err);
		return;
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}