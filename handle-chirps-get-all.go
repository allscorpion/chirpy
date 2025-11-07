package main

import (
	"net/http"
)


func (cfg *apiConfig) handleChirpsGetAll(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();

	chirps, err := cfg.dbQueries.GetAllChirps(req.Context());

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get chirps", err);
		return;
	}

	chirpsParsed := []ChirpParsed{};

	for _, chirp := range chirps {
		chirpsParsed = append(chirpsParsed, ChirpParsed{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		});
	}

	respondWithJSON(w, http.StatusOK, chirpsParsed)
}