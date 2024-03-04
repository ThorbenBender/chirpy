package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thorbenbender/chirpy/internal/auth"
	"github.com/thorbenbender/chirpy/internal/database"
)

func (cfg *apiConfig) HandlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetBearerToken(r.Header, "ApiKey")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No ApiKey found")
		return
	}
	if apiKey != cfg.ApiKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid Request")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldnt decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJson(w, http.StatusOK, struct{}{})
		return
	}

	err = cfg.DB.UpgradeUser(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Couldnt find user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldnt upgrade user")
		return
	}
	respondWithJson(w, http.StatusOK, struct{}{})
}
