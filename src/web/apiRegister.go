package web

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
)

func (a *API) register(c *fiber.Ctx) error {
	var (
		data database.User
		err  error
	)

	err = c.BodyParser(&data)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing data: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	for _, user := range util.USERS {
		if user == data.Name {
			return c.Status(fiber.StatusFound).JSON(fiber.Map{})
		}
	}

	err = database.NewUser(data, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error creating user: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	util.CacheUsers(a.DB)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
