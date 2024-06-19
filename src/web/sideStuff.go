package web

import (
	"fmt"
	"showmaster/src/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetLogLevel() {
	fiberLog.SetLevel(fiberLog.LevelInfo)
}

var (
	cfg    config.CFG = *config.GetConfig()
	server string     = fmt.Sprintf("%s:%d", cfg.Website.Host, cfg.Website.Port)

	APP *fiber.App = fiber.New(fiber.Config{
		ServerHeader: "showmaster:fiber",
		AppName:      "showmaster",
	}) // fiber app

	c = cors.New(cors.Config{
		AllowOrigins: strings.Join([]string{
			"*",
		}, ","),

		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
		}, ","),

		AllowHeaders: strings.Join([]string{
			"application/json",
		}, ","),

		AllowCredentials: false,
	}) // cors middleware

	mon = monitor.New(monitor.Config{
		Title: "ShowMaster Monitor",
	}) // Monitor

	auth = basicauth.New(basicauth.Config{
		Next: checkAuth,
		Users: map[string]string{
			"showmaster": "12345678",
		},
		Realm:           "Forbidden",
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}) // Basic Auth

	logs = logger.New(logger.Config{
		Next:     noLog,
		Format:   "[${ip}]:${port} ${status} - ${method} ${path}\n",
		TimeZone: "Europe/Berlin",
	})

	err error
)
