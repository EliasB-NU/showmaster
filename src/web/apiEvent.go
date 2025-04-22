package web

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
)

func (a *API) getEvents(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 1, "event", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	var events []database.Client
	if err := a.DB.Find(&events).Error; err != nil {
		log.Printf("Error getting events: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting events")
	}

	return c.Status(fiber.StatusOK).JSON(events)
}
