package web

import "github.com/gofiber/fiber/v2"

func getHealthcheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON("")
}
