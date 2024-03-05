package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/thorbenbender/chirpy/internal/auth"
	"github.com/thorbenbender/chirpy/internal/database"
)

// region -- handlerChirpRetrieve
func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	chirps := []database.Chirp{}
	if authorIDString == "" {
		dbChirps, err := cfg.DB.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldnt retrive chirps")
			return
		}
		chirps = dbChirps
	} else {

		authorID, err := strconv.Atoi(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldnt parse author id")
			return
		}
		dbChirps, err := cfg.DB.GetAuthorChirps(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldnt get chirps")
			return
		}
		chirps = dbChirps
	}
	orderBy := "asc"

	orderByParam := r.URL.Query().Get("sort")
	if orderByParam == "desc" {
		orderBy = "desc"
	}

	sort.Slice(chirps, func(i, j int) bool {
		if orderBy == "desc" {
			return chirps[i].ID > chirps[j].ID
		}
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJson(w, http.StatusOK, chirps)
}

// endregion -- handlerChirpRetrieve

// region -- handlerChirpsCreate
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	token, err := auth.GetBearerToken(r.Header, "Bearer")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "You are not logged in")
	}
	subject, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong jwt format")
		return
	}
	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt parse id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters")
		return
	}

	cleaned, err := validate_chirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := cfg.DB.CreateChirp(cleaned, userIDInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create chirp")
		return
	}
	respondWithJson(w, http.StatusCreated, chirp)
}

// endregion -- handlerChirpsCreate

// region -- handlerChirpRetrieve
func (cfg *apiConfig) handlerChirpRetrieve(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Id is in wrong format")
		return
	}
	chirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldnt retrieve chirp")
		return
	}

	respondWithJson(w, http.StatusOK, chirp)
}

// endregion -- handlerChirpRetrieve

// region -- handlerChirpDelete
func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldnt parse id")
		return
	}

	token, err := auth.GetBearerToken(r.Header, "Bearer")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT is in wrong format")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt validate JWT")
		return
	}
	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt parse user id")
		return
	}
	dbChirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldnt get chirp")
		return
	}

	if dbChirp.AuthorID != userIDInt {
		respondWithError(w, http.StatusForbidden, "You cant delete this chirp")
		return
	}

	err = cfg.DB.DeleteChirp(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt delete chirp")
		return
	}

	respondWithJson(w, http.StatusOK, struct{}{})
}

// endregion -- handlerChirpDelete
