package main

import (
	"encoding/json"
	"net/http"

	"github.com/allscorpion/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebook(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close();
	type EventData struct {
		UserId string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data EventData `json:"data"`
	}

	decoder := json.NewDecoder(req.Body);
	var params parameters;
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err);
		return;
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, struct{}{})
		return;
	}

	userId, err := uuid.Parse(params.Data.UserId);

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse user id", err);
		return;
	}

	_, err = cfg.dbQueries.GetUserById(req.Context(), userId);

	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to get user", err);
		return;
	}

	_, err = cfg.dbQueries.UpdateUserChirpyRed(req.Context(), database.UpdateUserChirpyRedParams{
		ID: userId,
		IsChirpyRed: true,
	});

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user", err);
		return;
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{});
}