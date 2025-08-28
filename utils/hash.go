package utils

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("failed to generate hash: %v", err)
		return "", errors.New("failed to generate hash")
	}
	return string(hash), nil
}

func CompareHashAndPassword(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Errorf("failed to compare hash and password with error: %v", err)
		return err
	}
	return nil
}