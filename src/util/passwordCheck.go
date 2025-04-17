package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashString hashes the password using bcrypt
func HashString(key string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckStringHash compares the plain password with the hashed one
func CheckStringHash(key, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(key))
	log.Printf("Error comparing password: %v\n", err)
	return err == nil
}
