package hash

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(pass string) (string, error) {
	const salt = 14
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), salt)
	if err != nil {
		log.Printf("ERR: %v\n", err)
		return "", err
	}

	return string(hash), nil
}

func ComparePasswordHash(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		log.Printf("ERR: %v\n", err)
		return err
	}

	return nil
}
