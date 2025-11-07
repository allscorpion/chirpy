package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, errMsg string, err error) {
	if err != nil {
		fmt.Printf("an error has occured %v\n", err);
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	errorResp := errorResponse{
		Error: errMsg,
	}

	respondWithJSON(w, statusCode, errorResp)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(statusCode);
	data, err := json.Marshal(payload);

	if err != nil {
		w.Write([]byte(fmt.Sprintf("unable to marshal json: %v", data)))
		return;
	}

	w.Write(data)
}