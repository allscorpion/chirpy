package main

import (
	"net/http"

	"github.com/google/uuid"
)


func (cfg *apiConfig) handleChirpsGet(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();

	chirpId, err := uuid.Parse(req.PathValue("chirpID"));

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirpID");
		return;
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpId);

	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to get chirps");
		return;
	}

	respondWithJSON(w, http.StatusOK, ChirpParsed{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}