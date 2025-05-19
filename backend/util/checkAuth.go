package util

import (
	"gorm.io/gorm"
	"log"
	"showmaster/backend/database"
	"strings"
)

// CheckAuth checks if the browser token is valid and deletes all, if there are multiple
func CheckAuth(headers map[string][]string, db *gorm.DB) bool {
	// Check if bearer token is present
	if headers["Authorization"] == nil {
		log.Println("No Authorization header")
		return false
	}
	// Check if bearer token is valid
	token := strings.TrimPrefix(headers["Authorization"][0], "Bearer ")
	if token == "" {
		log.Println("No token found")
		return false
	}

	// Check for an matching entry in the database
	var key []database.BrowserToken
	result := db.Where("token = ?", token).Find(&key)
	if result.Error != nil {
		log.Printf("Error getting token: %v\n", result.Error)
		return false
	}

	if len(key) > 1 {
		log.Println("More than one token found")
		go func() {
			for _, v := range key {
				err := db.Delete(v).Error
				if err != nil {
					log.Printf("Error deleting token: %v\n", err)
				}
			}
		}()
		return false
	}

	return true
}
