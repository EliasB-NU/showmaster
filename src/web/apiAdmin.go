package web

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
)

func (a *API) getUsers(c *fiber.Ctx) error {
	var (
		users []database.User
		err   error

		data []struct {
			Email       string `json:"email"`
			Name        string `json:"name"`
			Permissions struct {
				Login  bool `json:"login"`
				Admin  bool `json:"admin"`
				Events int  `json:"events"`
			}
		}
	)
	if !util.CheckPermissions(c.GetReqHeaders(), 3, util.Admin, a.DB) {
		return c.Status(fiber.StatusUnauthorized).JSON("unauthorized")
	}

	if err = a.DB.Find(&users).Error; err != nil {
		log.Printf("Error getting all users: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting all users")
	}

	for _, user := range users {
		var perms database.Permission
		if err = a.DB.Where("user_id = ?", user.ID).First(&perms).Error; err != nil {
			log.Printf("Error getting permissions for user %v: %v\n", user.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON("Error getting permissions for user")
		}

		data = append(data, struct {
			Email       string `json:"email"`
			Name        string `json:"name"`
			Permissions struct {
				Login  bool `json:"login"`
				Admin  bool `json:"admin"`
				Events int  `json:"events"`
			}
		}{
			Email: user.Email,
			Name:  user.Username,
			Permissions: struct {
				Login  bool `json:"login"`
				Admin  bool `json:"admin"`
				Events int  `json:"events"`
			}{
				Login:  perms.Login,
				Admin:  perms.Admin,
				Events: perms.Events,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

func (a *API) addUser(c *fiber.Ctx) error {
	var (
		data struct {
			Email       string `json:"email"`
			Name        string `json:"name"`
			Password    string `json:"password"`
			Permissions struct {
				Login  bool `json:"login"`
				Admin  bool `json:"admin"`
				Events int  `json:"events"`
			}
		}

		user  database.User
		perms database.Permission

		err error
	)
	if !util.CheckPermissions(c.GetReqHeaders(), 3, util.Admin, a.DB) {
		return c.Status(fiber.StatusUnauthorized).JSON("unauthorized")
	}
	// Parse and check body
	if err := c.BodyParser(&data); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON("Error parsing request body")
	}
	if data.Email == "" || data.Name == "" || data.Password == "" {
		log.Println("Error: Missing required fields")
		return c.Status(fiber.StatusBadRequest).JSON("Missing required fields")
	}
	if data.Permissions.Events < 0 || data.Permissions.Events > 3 {
		log.Println("Error: Invalid permissions")
		return c.Status(fiber.StatusBadRequest).JSON("Invalid permissions")
	}

	// Check if user already exists
	if err := a.DB.Where("email = ?", data.Email).First(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error checking if user exists: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON("Error checking if user exists")
		}
	} else {
		log.Println("Error: User already exists")
		return c.Status(fiber.StatusBadRequest).JSON("User already exists")
	}

	// Hash password
	hashedPassword, err := util.HashString(data.Password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error hashing password")
	}

	// Create user
	perms.Login = data.Permissions.Login
	perms.Admin = data.Permissions.Admin
	perms.Events = data.Permissions.Events
	user.Password = hashedPassword
	user.Email = data.Email
	user.Username = data.Name
	user.Perms = &perms
	if err = a.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error creating user")
	}

	return c.Status(fiber.StatusOK).JSON("")
}
