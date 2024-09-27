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
	"showmaster/src/util"
	"strings"
)

type API struct {
	DB      *sql.DB
	Clients map[*websocket.Conn]bool
}

var clients = make(map[*websocket.Conn]bool)

func InitWeb(db *sql.DB, cfg *config.CFG) {
	var (
		addr = fmt.Sprintf("%s:%d", cfg.Website.Host, cfg.Website.Port)

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
				"Access-Control-Allow-Headers",
			}, ","),

			AllowCredentials: false,
		})

		// Basic Auth
		auth = basicauth.New(basicauth.Config{
			Next:            checkAuth,
			Users:           util.GetAdminUsers(db),
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
	api := fiber.New()
	app.Mount("/api", api)
	a := API{
		DB:      db,
		Clients: clients,
	}
	// Websocket
	api.Get("/ws", websocket.New(a.WebsocketConnection))
	// Login/Register
	api.Post("/login", a.login)       // Login | <- incoming
	api.Post("/register", a.register) // Register | <- incoming
	// Admin
	api.Get("/admin/getusers", a.getUsers)              // Get all users | -> outgoing
	api.Patch("/admin/updateuser", a.updateUser)        // Update the permission level of a user | <- incoming
	api.Delete("/admin/deleteuser", a.deleteUser)       // Delete a user | <- incoming
	api.Delete("/admin/deleteproject", a.deleteProject) // Delete a project | <- incoming
	// Tables
	api.Get("/getprojects", a.getProjects)       // Get tables | -> outgoing
	api.Post("/newproject", a.newProject)        // New table | <- incoming
	api.Patch("/updateproject", a.updateProject) // Update table | <- incoming
	// Rows
	api.Get("/getrows:project", a.getRows)        // Get rows | -> outgoing
	api.Post("/newrow:project", a.newRow)         // New row | <- incoming
	api.Patch("/updaterow:project", a.updateRow)  // Update row | <- incoming
	api.Delete("/deleterow:project", a.deleteRow) // Delete row | <- incoming
	// Updates
	api.Get("/gethighlightedrow:project", a.getHighlightedRow)         // Get the currently highlighted row in a project | -> outgoing
	api.Patch("/updatehighlightedrow:project", a.updateHighlightedRow) // Update the currently highlighted row in a project | <- incoming
	api.Get("/gettimer:project", a.getTimer)                           // Get the current state and time of the timer | -> outgoing
	api.Patch("/updatetimer:project", a.updateTimer)                   // Update the timer status (update sent to clients via websocket) | <- incoming

	// Frontend
	app.Static("/", cfg.Website.Files)

	// Start fiber
	log.Println("Started ShowMaster V3")
	err = app.Listen(addr)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Web init error: %d\n", err)
	}
}
