package web

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
	"strconv"
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
		if a.LoadedEvents[u] {
			loadedId = u
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
	var (
		data = struct {
			Id          uint64    `json:"id"`
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
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid JSON")
	}
	if data.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Missing required fields")
	}

	// Find Event
	err = a.DB.First(&event, data.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("Event not found")
		}
		log.Printf("Error getting event: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting event")
	}

	event.Name = data.Name
	event.Description = data.Description
	event.Location = data.Location
	event.StartDate = data.StartDate
	event.EndDate = data.EndDate
	err = a.DB.Save(&event).Error
	if err != nil {
		log.Printf("Error updating event: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error updating event")
	}

	return c.Status(fiber.StatusOK).JSON("Event updated")
}

func (a *API) deleteEvent(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 3, "event", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	err := a.DB.Delete(&database.Event{}, c.Params("id")).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("Event not found")
		}
		log.Printf("Error deleting event: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error deleting event")
	}

	return c.Status(fiber.StatusOK).JSON("Event deleted")
}

func (a *API) activateEvent(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 3, "event", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	var id, _ = strconv.ParseUint(c.Params("id"), 10, 64)

	for e := range a.LoadedEvents {
		if a.LoadedEvents[e] {
			a.LoadedEvents[e] = false
		} else if e == id {
			a.LoadedEvents[e] = true
		}
	}

	return c.Status(fiber.StatusOK).JSON("Event activated")
}
