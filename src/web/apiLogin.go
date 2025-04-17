package web

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
	"unicode/utf8"
)

func (a *API) login(c *fiber.Ctx) error {
	var (
		data = struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			DeviceId string `json:"deviceId"`
		}{}

		err error
	)
	// Parse body && validate
	if err = c.BodyParser(&data); err != nil {
		log.Printf("Error parsing body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON("invalid request")
	}
	if data.Email == "" || data.Password == "" || utf8.RuneCountInString(data.DeviceId) != 16 {
		return c.Status(fiber.StatusBadRequest).JSON("invalid request")
	}

	// Check if user exists && password is correct
	var user database.User
	if err = a.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid credentials")
	}
	if !util.CheckStringHash(data.Password, user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON("invalid password")
	}

	// Check if user is already logged in
	var browserTokens []database.BrowserToken
	err = a.DB.Where("user_id = ?", user.ID).Find(&browserTokens).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error getting all browser tokens: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting all browser tokens")
	}
	if len(browserTokens) > 0 {
		var found bool
		for _, token := range browserTokens {
			if token.DeviceId == data.DeviceId {
				fmt.Printf("User %v is already logged in on device %v\n", user.ID, data.DeviceId)
				fmt.Println("All tokens will be deleted")
				// Delete all tokens for this user
				found = true
			}
		}
		if found {
			err = a.DB.Where("user_id = ? AND device_id = ?", user.ID, data.DeviceId).Delete(&database.BrowserToken{}).Error
			if err != nil {
				log.Printf("Error deleting browser tokens: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).JSON("Error deleting old browser tokens")
			}
		}
	}

	// Create new browser token
	var token = util.GenerateSessionToken()
	var hashedToken, _ = util.HashString(token)
	if token == "" {
		return c.Status(fiber.StatusInternalServerError).JSON("Error generating token")
	}
	browserToken := database.BrowserToken{
		DeviceId: data.DeviceId,
		Key:      hashedToken,
		User:     user,
		UserID:   user.ID,
	}
	err = a.DB.Create(&browserToken).Error
	if err != nil {
		log.Printf("Error creating browser token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error creating browser token")
	}

	// Get user permissions
	var perms database.Permission
	err = a.DB.Where("user_id = ?", user.ID).First(&perms).Error
	if err != nil {
		log.Printf("Error getting user permissions: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting user permissions")
	}

	// Check if user is allowed to login
	if !perms.Login {
		err = a.DB.Delete(&database.BrowserToken{UserID: user.ID, DeviceId: data.DeviceId, Key: token}).Error
		if err != nil {
			log.Printf("Error deleting browser token: %v\n", err)
		}
		return c.Status(fiber.StatusForbidden).JSON("user is not allowed to login")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
		"perms": fiber.Map{
			"login":  perms.Login,
			"admin":  perms.Admin,
			"events": perms.Events,
		},
	})
}

func (a *API) logout(c *fiber.Ctx) error {
	var (
		data = struct {
			DeviceId string `json:"deviceId"`
			Token    string `json:"token"`
		}{}

		err error
	)
	// Parse body && validate
	if err = c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid request")
	}
	if utf8.RuneCountInString(data.DeviceId) != 16 || data.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON("invalid request")
	}

	// Delete browser token
	var browserTokens []database.BrowserToken
	err = a.DB.Where("device_id = ?", data.DeviceId).Find(&browserTokens).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusForbidden).JSON("user is not logged in")
		}
		log.Printf("Error getting all browser tokens: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting all browser tokens")
	}
	for _, browserToken := range browserTokens {
		if util.CheckStringHash(data.Token, browserToken.Key) {
			err = a.DB.Delete(&browserToken).Error
			if err != nil {
				log.Printf("Error deleting browser token: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).JSON("Error deleting browser token")
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON("")
}

func (a *API) checkIfUserIsLoggedIn(c *fiber.Ctx) error {
	var (
		data = struct {
			DeviceId string `json:"deviceId"`
			Token    string `json:"token"`
		}{}

		err error
	)
	// Parse body && validate
	if err = c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid request")
	}
	if utf8.RuneCountInString(data.DeviceId) != 16 || data.Token == "" {
		log.Printf("Invalid request: %v\n", data)
		return c.Status(fiber.StatusBadRequest).JSON("invalid request")
	}

	// Check if user is logged in
	var browserTokens []database.BrowserToken
	err = a.DB.Where("device_id = ?", data.DeviceId).Find(&browserTokens).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusForbidden).JSON("user is not logged in")
		}
		log.Printf("Error getting all browser tokens: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting all browser tokens")
	}
	var found bool
	var userID uint64
	for _, browserToken := range browserTokens {
		if util.CheckStringHash(data.Token, browserToken.Key) {
			if found {
				err = a.DB.Where("device_id = ?", data.DeviceId).Delete(&database.BrowserToken{}).Error
				if err != nil {
					log.Printf("Error deleting browser token: %v\n", err)
					return c.Status(fiber.StatusInternalServerError).JSON("Error deleting browser token")
				}
				return c.Status(fiber.StatusForbidden).JSON("user is already logged in")
			} else {
				found = true
				userID = browserToken.UserID
			}
		}
	}
	if !found {
		return c.Status(fiber.StatusForbidden).JSON("user is not logged in")
	}

	// Get user permissions
	var perms database.Permission
	err = a.DB.Where("user_id = ?", userID).First(&perms).Error
	if err != nil {
		log.Printf("Error getting user permissions: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting user permissions")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"perms": fiber.Map{
			"login":  perms.Login,
			"admin":  perms.Admin,
			"events": perms.Events,
		},
	})
}
