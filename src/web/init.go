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

type API struct {
	DB      *sql.DB
	Clients map[*websocket.Conn]bool
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
				fiber.MethodPatch,
				fiber.MethodDelete,
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
	a := API{
		DB:      w.DB,
		Clients: clients,
	}
	// Login/Register
	app.Post("/api/login", a.login)       // Login | <- incoming
	app.Post("/api/register", a.register) // Register | <- incoming
	// Admin
	app.Get("/api/getusers", a.getUsers)              // Get all users | -> outgoing
	app.Patch("/api/updateuser", a.updateUser)        // Update the permission level of a user | <- incoming
	app.Delete("/api/deleteuser", a.deleteUser)       // Delete a user | <- incoming
	app.Delete("/api/deleteproject", a.deleteProject) // Delete a project | <- incoming
	// Tables
	app.Get("/api/gettables", a.getProjects)       // Get tables | -> outgoing
	app.Post("/api/newtable", a.newProject)        // New table | <- incoming
	app.Patch("/api/updatetable", a.updateProject) // Update table | <- incoming
	// Rows
	app.Get("/api/getrows", a.getRows)        // Get rows | -> outgoing
	app.Post("/api/newrow", a.newRow)         // New row | <- incoming
	app.Patch("/api/updaterow", a.updateRow)  // Update row | <- incoming
	app.Delete("/api/deleterow", a.deleteRow) // Delete row | <- incoming
	// Updates
	app.Get("/api/gethighlightedrow", a.getHighlightedRow)         // Get the currently highlighted row in a project | -> outgoing
	app.Patch("/api/updatehighlightedrow", a.updateHighlightedRow) // Update the currently highlighted row in a project | <- incoming
	app.Get("/api/gettimer", a.getTimer)                           // Get the current state and time of the timer | -> outgoing
	app.Patch("/api/updatetimer", a.updateTimer)                   // Update the timer status (update sent to clients via websocket) | <- incoming

	// Frontend
	app.Static("/", "./public/dist")

	// Start fiber
	err = app.Listen(addr)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Web init error: %d\n", err)
	}
}
