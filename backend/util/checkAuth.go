package util

import (
	"gorm.io/gorm"
	"log"
	"showmaster/backend/database"
	"strings"
)

// Used as the fars for permission check, so I don't have to remember the correct strings
var (
	ADMIN          = "admin"
	Director       = "director"
	FollowOperator = "follow_operator"
	Viewer         = "viewer"
)

// CheckAuth checks if the browser token is valid and deletes all, if there are multiple
func CheckAuth(headers map[string][]string, role string, level int, db *gorm.DB) bool {
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
	err := db.Where("token = ?", token).Find(&key).Error
	if err != nil {
		log.Printf("Error getting token: %v\n", err)
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

	// Get permissions
	var perm database.Permission
	err = db.Where("user_id = ?", key[0].UserID).First(&perm).Error
	if err != nil {
		log.Printf("Error getting permission: %v\n", err)
		return false
	}

	// Check the case for each role
	switch role {
	case ADMIN:
		if perm.Role == ADMIN {
			return true
		}
	case Director:
		if (perm.Role == Director && perm.Level == level) || perm.Role == ADMIN {
			return true
		}
	case FollowOperator:
		if (perm.Role == FollowOperator && perm.Level == level) || perm.Role == Director || perm.Role == ADMIN {
			return true
		}
	case Viewer:
		if (perm.Role == Viewer && perm.Level == level) || perm.Role == FollowOperator || perm.Role == Director || perm.Role == ADMIN {
			return true
		}
	default:
		log.Println("No role found")
		return false
	}

	return false
}
