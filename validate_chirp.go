package main

import (
	"encoding/json"
	"net/http"
	"strings"
)


func validate_chirp(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Body string `json:"body"`
  }
  type returnVals struct {
    CleanedBody string `json:"cleaned_body"`
  }
  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
    return
  }
  

  const maxChirpLength = 140
  if len(params.Body) > maxChirpLength {
    respondWithError(w, http.StatusBadRequest, "Chirp is too long")
    return
  }

  bannedWords := map[string]struct{}{
    "kerfuffle": {},
    "sharbert": {},
    "fornax": {},
  }

  cleaned := getCleanedBody(params.Body, bannedWords)
  

  respondWithJson(w, http.StatusOK, returnVals{
    CleanedBody: cleaned,
  }) 
}

func getCleanedBody(body string, bannedWords map[string]struct{}) string {
  words := strings.Split(body, " ")
  for i, word := range words {
    loweredWord := strings.ToLower(word)
    if _, ok := bannedWords[loweredWord]; ok {
      words[i] = "****"
    }
  }
  cleaned := strings.Join(words, " ")
  return cleaned
}
