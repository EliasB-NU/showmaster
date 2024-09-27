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
		i, p int
		err  error
	)

	if err := c.BodyParser(&data); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing data: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	i, p, err = database.CheckIfRegistered(data.Email, data.Password, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error checking if user is registered")
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	switch i {
	case 0:
		return c.Status(fiber.StatusForbidden).JSON("")
	case 1:
		return c.Status(fiber.StatusAccepted).JSON(p)
	case 2:
		return c.Status(fiber.StatusNotFound).JSON("")
	}

	return c.Status(fiber.StatusBadRequest).JSON("")
}
