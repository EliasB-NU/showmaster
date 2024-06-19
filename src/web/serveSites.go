package web

import (
	"fmt"
	"log"
	"showmaster/src/database"
	"showmaster/src/web/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/websocket/v2"
)

func StartTheWeb() {
	APP.Use(c)                                          // cors middleware
	APP.Use(auth)                                       // auth for monitor side
	APP.Use(logs)                                       // For logs
	APP.Use(healthcheck.New(healthcheck.ConfigDefault)) // Healthcheck
	APP.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}) // Websocket Middleware

	// Internal stuff
	APP.Get("/monitor", mon)

	// API Stuff
	APP.Get("/ws", websocket.New(api.WebsocketConnection))
	APP.Get("/api/getdata:project", api.GetData) // returns rows for specific table
	APP.Post("/api/updatehighlightedrow:project", api.UpdateHighlightedRow)

	APP.Get("/api/stopwatch-status:project", api.StopWatchStatus)  // Returns the current time & status
	APP.Post("/api/stopwatch-update:project", api.StopWatchUpdate) // Used for updating the stopwatch status

	APP.Get("/api/gettables", api.GetTables)          // returns all tables
	APP.Post("/api/newinsert:project", api.NewInsert) // For new insert into a table
	APP.Post("/api/newtable", api.NewTable)           // For new table

	// Sites
	APP.Static("/", "./public/projectSite") // The Project site
	allRowSites()                           // The sites for the individual projects

	// Init fiber
	err = APP.Listen(server)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error starting WebServer: %v\n", err)
	}
}

func allRowSites() {
	var tables []string = database.GetTables()
	for _, table := range tables {
		var URL string = fmt.Sprintf("/%s", table)
		APP.Static(URL, "./public/rowSite")
	}
}

func CheckForNewRowSite() {
	var oldtables []string = database.GetTables()
	for {
		var newTables []string = database.GetTables()
		if len(newTables) > len(oldtables) {
			for _, table1 := range newTables {
				for _, table2 := range oldtables {
					if table1 == table2 {
						continue
					} else {
						oldtables = newTables
						var URL string = fmt.Sprintf("/%s", table1)
						APP.Static(URL, "./public/rowSite")
					}
				}
			}
		}
	}
}
