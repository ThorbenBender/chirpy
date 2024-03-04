package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "chirpy-access"
	TokenTypeRefresh TokenType = "chirpy-refresh"
)

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(
	userID int,
	tokenSecret string,
	expiresIn time.Duration,
	tokenType TokenType,
) (string, error) {
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(tokenType),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprintf("%d", userID),
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return "", err
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeAccess) {
		return "", errors.New("invalid issuer")
	}
	return userIDString, nil
}

func RefreshToken(tokenString, tokenSecret string) (string, error) {
	claimStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimStruct,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", nil
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeRefresh) {
		return "", errors.New("invalid issuer")
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return "", err
	}

	newToken, err := MakeJWT(userID, tokenSecret, time.Hour, TokenTypeAccess)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

func GetBearerToken(headers http.Header, authToken string) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No auth header included")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != authToken {
		return "", errors.New("Malformed authorization header")
	}
	return splitAuth[1], nil
}
