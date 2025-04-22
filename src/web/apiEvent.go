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

	var events []database.Event
	if err := a.DB.Find(&events).Error; err != nil {
		log.Printf("Error getting events: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting events")
	}

	var loadedId uint64 = 0
	for u := range a.LoadedEvents {
		if a.LoadedEvents[u] == true {
			loadedId = u
			break
		}
		continue
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"events": events,
		"active": loadedId,
	})
}

func (a *API) createEvent(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON("W.I.P.")
}

func (a *API) updateEvent(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON("W.I.P.")
}

func (a *API) deleteEvent(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON("W.I.P.")
}

func (a *API) activateEvent(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON("W.I.P.")
}
