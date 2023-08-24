package main

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashStringBcrypt(data string) (string, error) {
	hashedString, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash string: %w", err)
	}
	return string(hashedString), nil
}

func VerifyStringBcrypt(hashedData string, candidateData string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(candidateData))
}

func HashStringSha256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return string(hash[:])
	//h := sha256.New()
	//h.Write([]byte(data))
	//hash := h.Sum(nil)
	//return string(hash)
}
