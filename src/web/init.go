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

	LoadedEvents map[uint64]bool

	Stopwatch *util.Stopwatch
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

		LoadedEvents: make(map[uint64]bool),

		Stopwatch: util.NewStopwatch(),
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
			a.LoadedEvents[event.ID] = false
		}
	}
	// Websocket
	api.Get("/ws", websocket.New(a.WebsocketConnection))
	// Login
	api.Post("/login", a.login)                      // <- Email&Password || ->returns new session token
	api.Delete("/logout", a.logout)                  // <- Token, deletes session
	api.Post("/checkLogin", a.checkIfUserIsLoggedIn) // -> Bool&Perms, checks if the session is valid and returns the users permissions
	// Admin
	api.Get("/users", a.getUsers)                      // <- Token || -> returns all users
	api.Post("/users/create", a.addUser)               // <- Token&Data, creates a new user
	api.Post("/users/update", a.updateUser)            // <- Token&Data updates a user
	api.Delete("/users/delete/:id", a.deleteUser)      // <- Token&Id, deletes a user
	api.Get("/clients", a.getClients)                  // <- Token || -> returns all clients
	api.Post("/clients/create", a.createClient)        // <- Token&Data, creates a new client
	api.Post("/clients/update", a.updateClient)        // <- Token&Id&Data, updates a client
	api.Delete("/clients/delete/:id", a.deleteClient)  // Token&Id, deletes client
	api.Get("/clients/newToken/:id", a.createNewToken) // <- Token&Id || -> generates a new auth token and deletes the old one
	api.Get("/clients/byType/:type", a.getClientType)  // <- Token&Type || -> returns all clients with a specific type
	// Events
	api.Get("/events", a.getEvents)                  // <- Token || -> All events & the currently loaded event
	api.Post("/event/create", a.createEvent)         // <- Token&Data, creates new event
	api.Post("/event/update", a.updateClient)        // <- Token&Data, updates an event
	api.Delete("/event/delete/:id", a.deleteEvent)   // <- Token&Id, deletes an event
	api.Post("/event/activate/:id", a.activateEvent) // <- Token&Id, activates an event
	// Scenes
	api.Get("/scenes/:id", a.getScenes)                    // <- Token || -> All scenes from an event
	api.Post("/scenes/create/:id", a.createScene)          // <- Token&Data&id, creates a new scene in an event
	api.Post("/scenes/update/:id", a.updateScene)          // <- Token&Data&id, updates a scene in an event
	api.Post("/scenes/delete/:id/:sceneId", a.deleteScene) // <- Token&id&sceneId, deletes a scene by its id from an event
	// Timer
	api.Get("/timer", a.getTimer)      // <- Token || -> time.Duration of the timer of the active project
	api.Post("timer", a.updateTimer)   // <-Token&Command, updates the state of the timer of the active project (start, stop, resume)
	api.Delete("timer", a.deleteTimer) // <- Token, deletes the time.Duration of the timer of the active project (extra permission)
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
