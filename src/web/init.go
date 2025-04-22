package web

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
	"log"
	"showmaster/src/config"
	"showmaster/src/database"
	"showmaster/src/util"
	"strings"
)

type API struct {
	DB      *gorm.DB
	CFG     *config.Config
	Clients map[*websocket.Conn]bool

	LoadedEvents map[string]bool
}

func InitWeb(cfg *config.Config, db *gorm.DB, mst *util.MST) {
	var (
		addrShowMaster = "0.0.0.0:3000"

		err error

		// Fiber
		showMasterApp = fiber.New(fiber.Config{
			ServerHeader: "showmaster:fiber",
			AppName:      "showmaster",
		})

		// Cors
		c = cors.New(cors.Config{
			AllowOrigins: strings.Join([]string{
				"*",
			}, ","),

			AllowHeaders: strings.Join([]string{
				"Origin",
				"Content-Type",
				"Accept",
			}, ","),

			AllowMethods: strings.Join([]string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodPatch,
				fiber.MethodDelete,
			}, ","),

			AllowCredentials: false,
		})

		// Monitor
		mon = monitor.New(monitor.Config{
			Title: "Showmaster Monitor",
		})
	)
	// Internal tools
	showMasterApp.Use(c)                                          // Cors Middleware
	showMasterApp.Use(healthcheck.New(healthcheck.ConfigDefault)) // Healthcheck Middleware
	showMasterApp.Use("/api/ws", func(c *fiber.Ctx) error {       // Websocket Middleware
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	showMasterApp.Use("/monitor", mon)                // Monitor
	showMasterApp.Get("/healthcheck", getHealthcheck) // Healthcheck

	// API
	api := fiber.New()
	showMasterApp.Mount("/api", api)
	a := API{
		DB:      db,
		CFG:     cfg,
		Clients: make(map[*websocket.Conn]bool),

		LoadedEvents: make(map[string]bool),
	}
	// Get all events
	var events []database.Event
	err = db.Find(&events).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("No events found")
		}
		log.Fatal("Error getting events: ", err)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// This code is for loading the currently selected scene in the redis cache
		for _, event := range events {
			a.LoadedEvents[event.Name] = false
		}
	}
	// Websocket
	api.Get("/ws", websocket.New(a.WebsocketConnection))
	// Login
	api.Post("/login", a.login)                      // <- Email&Password, returns new session token
	api.Delete("/logout", a.logout)                  // <- Token, deletes session
	api.Post("/checkLogin", a.checkIfUserIsLoggedIn) // -> Bool&Perms, checks if the session is valid and returns the users permissions
	// Admin
	api.Get("/getUsers", a.getUsers)                 // <- Token, returns all users
	api.Post("/createUser", a.addUser)               // <- Token&Email&Name&Password, creates a new user
	api.Post("/updateUser", a.updateUser)            // <- Token&Id&Email&Name&Password, updates a user
	api.Delete("/deleteUser/:id", a.deleteUser)      // <- Token&Id, deletes a user
	api.Get("/getClients", a.getClients)             // <- Token, returns all clients
	api.Post("/createClient", a.createClient)        // <- Token&Name&Type&IP&Port, creates a new client
	api.Post("/updateClient", a.updateClient)        // <- Token&Id&Name&Type&IP&Port, updates a client
	api.Delete("/deleteClient/:id", a.deleteClient)  // Token&Id, deletes client
	api.Get("/clientNewToken/:id", a.createNewToken) // <- Token&Id, generates a new auth token and deletes the old one
	api.Get("/getClientType/:type", a.getClientType) // <- Token&Type, returns all clients with a specific type
	// Events

	// Scenes

	// Active for low latency
	// Web
	showMasterApp.Static("/", "./web/dist")

	mst.ElapsedTime()
	// Start server
	log.Println("Started Showmaster V3")
	err = showMasterApp.Listen(addrShowMaster)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
