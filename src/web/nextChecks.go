package web

import (
	"github.com/gofiber/fiber/v2"
)

func checkAuth(c *fiber.Ctx) bool {
	toCheck := [6]string{"/monitor", "/api/admin/getusers", "/api/admin/updateuser", "/api/admin/deleteuser", "/api/admin/deleteproject"}
	for _, v := range toCheck {
		if c.OriginalURL() == v {
			return false
		}
	}
	return true
}

func noLog(c *fiber.Ctx) bool {
	return c.OriginalURL() == "/monitor"
}
