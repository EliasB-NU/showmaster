package web

import (
	"github.com/gofiber/fiber/v2"
	"showmaster/src/util"
)

func (a *API) getTimer(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 1, "event", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("Invalid Permission")
	}

	return c.Status(fiber.StatusOK).JSON(a.Stopwatch.ElapsedSeconds())
}

func (a *API) updateTimer(c *fiber.Ctx) error {
	var (
		data = struct {
			Status string `json:"status"`
		}{}
	)
	if !util.CheckPermissions(c.GetReqHeaders(), 1, "event", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("Invalid Permission")
	}
	// Parse&validate body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid Request")
	}
	if data.Status != "start" && data.Status != "stop" && data.Status != "resume" {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid Request")
	}

	switch data.Status {
	case "start":
		a.Stopwatch.Start()
		a.SendMessage([]byte(`{"stopwatch":"start"}`))
		break
	case "stop":
		a.Stopwatch.Stop()
		a.SendMessage([]byte(`{"stopwatch":"stop"}`))
		break
	case "resume":
		a.Stopwatch.Resume()
		a.SendMessage([]byte(`{"stopwatch":"resume"}`))
		break
	default:
		return c.Status(fiber.StatusBadRequest).JSON("Invalid Request")
	}

	return c.Status(fiber.StatusOK).JSON("State changed")
}

func (a *API) deleteTimer(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 2, "event", a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("Invalid Permission")
	}

	a.Stopwatch.Reset()
	a.SendMessage([]byte(`{"stopwatch":"reset"}`))

	return c.Status(fiber.StatusOK).JSON("Stopwatch reset")
}
