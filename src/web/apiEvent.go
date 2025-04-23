package web

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
	"time"
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
	var (
		data = struct {
			Name        string    `json:"name"`
			Description string    `json:"description"`
			Location    string    `json:"location"`
			StartDate   time.Time `json:"startDate"`
			EndDate     time.Time `json:"endDate"`
		}{}

		event database.Event
		err   error
	)
	if !util.CheckPermissions(c.GetReqHeaders(), 2, "event", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}
	// Parse&validate data
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid JSON")
	}
	if data.Name == "" || data.StartDate.IsZero() && data.EndDate.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON("Missing required fields")
	}

	event.Name = data.Name
	event.Description = data.Description
	event.Location = data.Location
	event.StartDate = data.StartDate
	event.EndDate = data.EndDate

	err = a.DB.Create(&event).Error
	if err != nil {
		log.Printf("Error creating event: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error creating event")
	}

	return c.Status(fiber.StatusOK).JSON("Event created")
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
