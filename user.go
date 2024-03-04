package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/thorbenbender/chirpy/internal/auth"
)

type User struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type authResponse struct {
	User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	encryptedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, encryptedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}
	respondWithJson(w, http.StatusCreated, User{
		Email:       user.Email,
		ID:          user.ID,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters")
	}
	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}
	err = auth.CheckPassword(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong password")
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour, auth.TokenTypeAccess)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create access JWT")
	}

	refreshToken, err := auth.MakeJWT(
		user.ID,
		cfg.JWTSecret,
		time.Hour*24*30*6,
		auth.TokenTypeRefresh,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create refresh JWT")
	}

	respondWithJson(w, http.StatusOK, authResponse{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	token, err := auth.GetBearerToken(r.Header, "Bearer")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt find jwt")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt validate JWT")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters")
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt hash password")
		return
	}
	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt parse id")
		return
	}
	user, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt update user")
		return
	}
	respondWithJson(w, http.StatusOK, User{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
