package main

import (
	"errors"
	"strings"
)


func validate_chirp(body string) (string, error) {
    

  const maxChirpLength = 140
  if len(body) > maxChirpLength {
    return "", errors.New("Chirp is too long")
  }

  bannedWords := map[string]struct{}{
    "kerfuffle": {},
    "sharbert": {},
    "fornax": {},
  }

  cleaned := getCleanedBody(body, bannedWords)
  
  return cleaned, nil 
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
