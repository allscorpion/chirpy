package main

import (
	"net/http"

	"github.com/allscorpion/chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handleChirpsGetAll(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();

	author_id := req.URL.Query().Get("author_id");

	var (
        chirps []database.Chirp
        err    error
    )

    if author_id == "" {
        chirps, err = cfg.dbQueries.GetAllChirps(req.Context())
    } else {
        var userID uuid.UUID
        userID, err = uuid.Parse(author_id)
        if err != nil {
            respondWithError(w, http.StatusBadRequest, "invalid author_id", err)
            return
        }
        chirps, err = cfg.dbQueries.GetAllChirpsForUser(req.Context(), userID)
    }

	if err != nil {
        respondWithError(w, http.StatusInternalServerError, "failed to get chirps", err)
        return
    }

	chirpsParsed := make([]ChirpParsed, 0, len(chirps))

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