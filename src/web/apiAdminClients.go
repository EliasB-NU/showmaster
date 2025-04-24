package web

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
)

func (a *API) getClients(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 3, util.Admin, a.DB) {
		return fiber.NewError(fiber.StatusForbidden)
	}

	var clients []database.Client
	err := a.DB.Find(&clients).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("")
		}
		log.Printf("Failed to get clients: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting clients")
	}

	return c.Status(fiber.StatusOK).JSON(clients)
}

func (a *API) getClientType(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 2, util.Events, a.DB) {
		return fiber.NewError(fiber.StatusForbidden)
	}

	var clients []database.Client
	err := a.DB.Where("type = ?", c.Params("type")).Find(&clients).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("")
		}
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
	if !util.CheckPermissions(c.GetReqHeaders(), 3, util.Admin, a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}
	// Parse & validate body
	if err = c.BodyParser(&data); err != nil {
		log.Printf("Failed to parse body: %v\n", err)
	}
	if data.Name == "" || data.IP == "" || data.Port < 1 || data.Port > 65535 {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid data")
	}
	if data.Type != "osc" && data.Type != "midi" && data.Type != "intracast" && data.Type != "gpio" {
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
	var (
		data = struct {
			Id   int `json:"id"`
			Name string
			Type string
			IP   string
			Port int
		}{}

		client database.Client
		err    error
	)
	if !util.CheckPermissions(c.GetReqHeaders(), 3, util.Admin, a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}
	// Parse & validate body
	if err := c.BodyParser(&data); err != nil {
		log.Printf("Failed to parse body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON("Invalid data")
	}
	if data.Id == 0 || data.Name == "" || data.IP == "" || data.Port < 1 || data.Port > 65535 {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid data")
	}
	if data.Type != "osc" && data.Type != "midi" && data.Type != "intracast" {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid data type")
	}

	err = a.DB.First(&client, data.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("")
		}
		log.Printf("Failed to get client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting client")
	}

	client.Name = data.Name
	client.Type = data.Type
	client.IP = data.IP
	client.Port = data.Port

	err = a.DB.Save(&client).Error
	if err != nil {
		log.Printf("Failed to update client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error updating client")
	}

	return c.Status(fiber.StatusOK).JSON("")
}

func (a *API) deleteClient(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 3, "admin", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	err := a.DB.Delete(&database.Client{}, c.Params("id")).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("No client found")
		}
		log.Printf("Failed to delete client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error deleting client")
	}

	return c.Status(fiber.StatusOK).JSON("Client deleted")
}

func (a *API) createNewToken(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 3, "admin", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	var newToken = util.GenerateSessionToken()
	var client database.Client
	err := a.DB.Where("id = ?", c.Params("id")).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("No client found")
		}
		log.Printf("Failed to find client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error creating client")
	}

	client.Token = newToken
	err = a.DB.Save(&client).Error
	if err != nil {
		log.Printf("Failed to update client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error creating client")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": client.Token,
	})
}
