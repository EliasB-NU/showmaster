package web

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"showmaster/src/database"
)

type userLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// login Password match: 202; Password mismatch: 403; Unregistered: 404
func (a *API) login(c *fiber.Ctx) error {
	var (
		data userLoginData
		i    int
		err  error
	)

	if err := c.BodyParser(&data); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing data: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	i, err = database.CheckIfRegistered(data.Email, data.Password, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error checking if user is registered")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	switch i {
	case 0:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{})
	case 1:
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{})
	case 2:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
}
