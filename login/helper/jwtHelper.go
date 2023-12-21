package helper

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

func CreateJWT(id string) (string, error) {
	signingKey := uuid.NewString()
	aToken := jwt.New(jwt.SigningMethodHS256)
	claims := aToken.Claims.(jwt.MapClaims)
	claims["token"] = signingKey
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["info"] = "wonju"

	tk, err := aToken.SignedString([]byte("SECRETKEY"))
	if err != nil {
		return "", err
	}
	return tk, nil
}
