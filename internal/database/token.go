package database

import "time"

type Revocation struct {
	Token     string    `json:"token"`
	RevokedAt time.Time `json:"revoked_at"`
}

func (db *DB) IsTokenRevoked(token string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	revocation, ok := dbStructure.Revocations[token]
	if !ok {
		return false, nil
	}
	if revocation.RevokedAt.IsZero() {
		return false, nil
	}
	return true, nil
}

func (db *DB) RevokeToken(tok string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	revocation := Revocation{
		Token:     tok,
		RevokedAt: time.Now().UTC(),
	}
	dbStructure.Revocations[tok] = revocation

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
