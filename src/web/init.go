package web

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/websocket/v2"
	"log"
	"showmaster/src/config"
	"strings"
)

type WEB struct {
	DB  *sql.DB
	CFG *config.CFG
}

var clients = make(map[*websocket.Conn]bool)

func (w *WEB) InitWeb() {
	var (
		addr = fmt.Sprintf("%s:%d", w.CFG.Website.Host, w.CFG.Website.Port)

		err error

		// Fiber app
		app = fiber.New(fiber.Config{
			ServerHeader: "showmaster:fiber",
			AppName:      "showmaster",
		})

		// Cors middleware
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
		})

		// Basic Auth
		auth = basicauth.New(basicauth.Config{
			Next: checkAuth,
			Users: map[string]string{
				"showmaster": "12345678",
			},
			Realm:           "Forbidden",
			ContextUsername: "_user",
			ContextPassword: "_pass",
		})

		// Logger
		logs = logger.New(logger.Config{
			Next:     noLog,
			Format:   "[${ip}]:${port} ${status} - ${method} ${path}\n",
			TimeZone: "Europe/Berlin",
		})

		// Monitor
		mon = monitor.New(monitor.Config{
			Title: "ShowMaster Monitor",
		})
	)

	// Internal tools
	app.Use(c)                                          // Cors middleware
	app.Use(auth)                                       // Basic auth for monitor side
	app.Use(logs)                                       // Logger
	app.Use(healthcheck.New(healthcheck.ConfigDefault)) // Healthcheck
	app.Use("/ws", func(c *fiber.Ctx) error {           // Websocket middleware
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/monitor", mon) // Monitor

	// API

	// Frontend
	app.Static("/", "./public/dist")

	// Start fiber
	err = app.Listen(addr)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Web init error: %d\n", err)
	}
}
