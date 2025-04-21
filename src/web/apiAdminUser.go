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
			Id          uint64 `json:"id"`
			Email       string `json:"email"`
			Name        string `json:"name"`
			Permissions struct {
				Login  bool `json:"login"`
				Admin  bool `json:"admin"`
				Events int  `json:"events"`
			} `json:"permissions"`
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
			Id          uint64 `json:"id"`
			Email       string `json:"email"`
			Name        string `json:"name"`
			Permissions struct {
				Login  bool `json:"login"`
				Admin  bool `json:"admin"`
				Events int  `json:"events"`
			} `json:"permissions"`
		}{
			Id:    user.ID,
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
			} `json:"permissions"`
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

	// Hash password
	hashedPassword := util.HashString(data.Password)

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

func (a *API) deleteUser(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 3, util.Admin, a.DB) {
		return c.Status(fiber.StatusUnauthorized).JSON("unauthorized")
	}

	// Delete user
	if err := a.DB.Delete(&database.User{}, c.Params("id")).Error; err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error deleting user")
	}
	// Delete permission
	if err := a.DB.Where("user_id = ?", c.Params("id")).Delete(&database.Permission{}).Error; err != nil {
		log.Printf("Error deleting permissions: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error deleting permissions")
	}
	// Delete browser tokens
	if err := a.DB.Where("user_id = ?", c.Params("id")).Delete(&database.BrowserToken{}).Error; err != nil {
		log.Printf("Error deleting browser token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error deleting browser token")
	}

	return c.Status(fiber.StatusOK).JSON("")
}

func (a *API) updateUser(c *fiber.Ctx) error {
	var (
		data struct {
			Id          uint64 `json:"id"`
			Email       string `json:"email"`
			Name        string `json:"name"`
			Password    string `json:"password"`
			Permissions struct {
				Login  bool `json:"login"`
				Admin  bool `json:"admin"`
				Events int  `json:"events"`
			} `json:"permissions"`
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

	// Update user
	if err = a.DB.Where("id = ?", data.Id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error: User not found")
			return c.Status(fiber.StatusNotFound).JSON("User not found")
		}
		log.Printf("Error getting user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting user")
	}

	user.Email = data.Email
	user.Username = data.Name
	if data.Password != "" {
		user.Password = util.HashString(data.Password)
	}
	if err = a.DB.Save(&user).Error; err != nil {
		log.Printf("Error updating user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error updating user")
	}

	// Update permissions
	if err = a.DB.Where("user_id = ?", data.Id).First(&perms).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error: Permissions not found")
			return c.Status(fiber.StatusNotFound).JSON("Permissions not found")
		}
		log.Printf("Error getting permissions: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting permissions")
	}
	perms.Login = data.Permissions.Login
	perms.Admin = data.Permissions.Admin
	perms.Events = data.Permissions.Events
	if err = a.DB.Save(&perms).Error; err != nil {
		log.Printf("Error updating permissions: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error updating permissions")
	}

	return c.Status(fiber.StatusOK).JSON("")
}
