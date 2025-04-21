package web

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
)

func (a *API) getClients(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 3, "admin", a.DB) {
		return fiber.NewError(fiber.StatusForbidden)
	}

	var clients []database.Client
	if err := a.DB.Find(&clients); err != nil {
		log.Printf("Failed to get clients: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting clients")
	}

	return c.Status(fiber.StatusOK).JSON(clients)
}

func (a *API) createClient(c *fiber.Ctx) error {
	var (
		data = struct {
			Name string `json:"name"`
			Type string `json:"type"`
			IP   string `json:"ip"`
			Port int    `json:"port"`
		}{}

		client database.Client

		err error
	)
	if !util.CheckPermissions(c.GetReqHeaders(), 3, "admin", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}
	// Parse & validate body
	if err = c.BodyParser(&data); err != nil {
		log.Printf("Failed to parse body: %v\n", err)
	}
	if data.Name == "" || data.IP == "" || data.Port < 1 || data.Port > 65535 {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid data")
	}
	if data.Type != "osc" && data.Type != "midi" && data.Type != "llls" && data.Type != "gpio" {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid data type")
	}

	client.Name = data.Name
	client.Type = data.Type
	client.IP = data.IP
	client.Port = data.Port
	client.Token = util.GenerateSessionToken()
	if err = a.DB.Create(&client).Error; err != nil {
		log.Printf("Failed to create client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error creating client")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": client.Token,
	})
}

func (a *API) updateClient(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON("")
}

func (a *API) deleteClient(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON("")
}

func (a *API) createNewToken(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON("")
}
