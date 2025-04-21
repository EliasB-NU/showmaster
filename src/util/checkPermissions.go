package util

import (
	"gorm.io/gorm"
	"log"
	"showmaster/src/database"
	"strings"
)

// CheckPermissions checks if the user has the required permissions
func CheckPermissions(headers map[string][]string, level int, subPart string, db *gorm.DB) bool {
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
	var browserToken database.BrowserToken
	result := db.Where("token = ?", token).Find(&browserToken)
	if result.Error != nil {
		log.Printf("Error getting browser tokens: %v\n", result.Error)
		return false
	}

	// Get the permissions for the user
	var perms database.Permission
	result = db.Where("user_id = ?", browserToken.UserID).First(&perms)
	if result.Error != nil {
		log.Printf("Error getting permissions: %v\n", result.Error)
		return false
	}

	// Check if the user is an admin
	if perms.Admin {
		return true
	}

	// Check if the user has the required permissions
	switch subPart {
	case Events:
		if perms.Events >= level || perms.Admin {
			return true
		} else {
			return false
		}
	case Admin:
		if perms.Admin {
			return true
		} else {
			return false
		}
	case Login:
		if perms.Login {
			return true
		} else {
			return false
		}
	default:
		if perms.Admin {
			return true
		} else {
			return false
		}
	}
}
