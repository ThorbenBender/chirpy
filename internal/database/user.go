package database

import (
	"errors"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
		ID:       id,
		Email:    email,
		Password: password,
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
