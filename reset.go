package main

import "net/http"

func (cfg *apiConfig) reset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden);
		w.Write([]byte("invalid environment"))
		return;
	}
	 
	cfg.fileserverHits.Store(0);
	err := cfg.dbQueries.DeleteUsers(req.Context());

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError);
		w.Write([]byte("failed to delete users"));
		return;
	}

	w.WriteHeader(http.StatusOK);
	w.Write([]byte("Hits reset to 0 and database reset"))
}