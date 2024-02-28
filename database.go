package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)


type DB struct {
  path string
  mux *sync.RWMutex
}

type DBStructure struct {
  Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
  ID int `json:"id"`
  Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
  db := DBStructure{
    Chirps: map[int]Chirp{},
  }
  data, err := json.Marshal(db)
  if err != nil {
    return nil, err
  }

  file := os.WriteFile("./data/database.json", data, 0666)
  fmt.Println(file)
  return &DB{
    path: path,
    mux: &sync.RWMutex{},
  }, nil
}
