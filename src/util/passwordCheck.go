package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashString hashes the password using bcrypt
func HashString(key string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing String: %v\n", err)
		return ""
	}
	return string(bytes)
}

// CheckStringHash compares the plain password with the hashed one
func CheckStringHash(key, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(key))
	if err != nil {
		log.Printf("Error comparing StringHash: %v\n", err)
		return false
	}
	return true
}
