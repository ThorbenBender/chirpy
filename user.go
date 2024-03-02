package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/thorbenbender/chirpy/internal/auth"
)

type authResponse struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	Token string `json:"token"`
}
type userResponse struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
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
	respondWithJson(w, http.StatusCreated, userResponse{
		Email: user.Email,
		ID:    user.ID,
	})
}

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:   "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(
			time.Now().UTC().Add(time.Duration(params.ExpiresInSeconds) * time.Second),
		),
		Subject: strconv.Itoa(user.ID),
	})
	authToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt craete auth token")
		return
	}
	respondWithJson(w, http.StatusOK, authResponse{
		Token: authToken,
		ID:    user.ID,
		Email: user.Email,
	})
}

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt find jwt")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.JWTSecret)
	fmt.Println(err)
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
	respondWithJson(w, http.StatusOK, userResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}
