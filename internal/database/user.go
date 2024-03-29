package database

import (
	"errors"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

var ErrAlreadyExists = errors.New("User already exists")

func (db *DB) DoesUserExist(email string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, nil
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (db *DB) CreateUser(email, password string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		ID:          id,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
	}
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, nil
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrNotExist
}

func (db *DB) UpdateUser(userID int, email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStructure.Users[userID]
	if !ok {
		return User{}, ErrNotExist
	}
	user.Email = email
	user.Password = password
	dbStructure.Users[userID] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) UpgradeUser(userID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	user, ok := dbStructure.Users[userID]
	if !ok {
		return ErrNotExist
	}
	user.IsChirpyRed = true
	dbStructure.Users[userID] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
