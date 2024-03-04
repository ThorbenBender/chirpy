package main

import (
	"net/http"

	"github.com/thorbenbender/chirpy/internal/auth"
)

type refreshResponse struct {
	Token string `json:"token"`
}
type revokeResponse struct {
	Revoked bool `json:"revoked"`
}

func (cfg *apiConfig) HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header, "Bearer")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt find refresh token")
		return
	}

	isRevoked, err := cfg.DB.IsTokenRevoked(refreshToken)

	if isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is revoked")
		return
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt validate JWT")
		return
	}
	respondWithJson(w, http.StatusOK, refreshResponse{
		Token: accessToken,
	})
}

func (cfg *apiConfig) HandleTokenRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header, "Bearer")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt find refresh token")
		return
	}
	err = cfg.DB.RevokeToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt revoke token")
		return
	}

	respondWithJson(w, http.StatusOK, struct{}{})
}
