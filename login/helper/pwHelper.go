package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(userPw, hashVal string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPw), []byte(hashVal))
	if err != nil {
		return false
	} else {
		return true
	}
}
