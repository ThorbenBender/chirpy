package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thorbenbender/chirpy/internal/database"
)

type apiConfig struct {
  fileServerHits int
  DB *database.DB
}



func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Content-Type", "text/html; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(fmt.Sprintf(`
  <html>
    <body>
      <h1>Welcome, Chirpy Admin</h1>
      <p>Chirpy has been visited %d times!</p>
    </body>
  </html>
  `,cfg.fileServerHits)))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
  cfg.fileServerHits = 0
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    cfg.fileServerHits += 1
    next.ServeHTTP(w, r)
  })
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
  chirps, err := cfg.DB.GetChirps()
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldnt retrive chirps")
  }
  
  sort.Slice(chirps, func(i, j int) bool {
    return chirps[i].ID < chirps[j].ID
  })
  respondWithJson(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Body string `json:"body"`
  }
  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters")
    return
  }

  cleaned, err := validate_chirp(params.Body)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }
  chirp, err := cfg.DB.CreateChirp(cleaned)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldnt create chirp")
    return
  }
  respondWithJson(w, http.StatusCreated, chirp)
}

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
