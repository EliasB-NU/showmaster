package web

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"showmaster/backend/config"
	"showmaster/backend/util"
	"strings"
)

type API struct {
	PSQL *gorm.DB
	RDB  *redis.Client

	Clients map[*websocket.Conn]bool

	CTX context.Context
	CFG *config.Config
}

func InitWeb(psql *gorm.DB, rdb *redis.Client, ctx context.Context, cfg *config.Config, mst *util.MST) {
	var (
		addr = "0.0.0.0:3000"

		err error

		showmasterAPP = fiber.New(fiber.Config{
			ServerHeader: "showmaster:fiber",
			AppName:      "showmaster",
		})

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

		a = API{
			PSQL:    psql,
			RDB:     rdb,
			Clients: make(map[*websocket.Conn]bool),

			CTX: ctx,
			CFG: cfg,
		}
	)

	// Internal tools
	showmasterAPP.Use(c)                                          // Cors Middleware
	showmasterAPP.Use(healthcheck.New(healthcheck.ConfigDefault)) // Healthcheck Middleware
	showmasterAPP.Use("/api/ws", func(c *fiber.Ctx) error {       // Websocket Middleware
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	showmasterAPP.Get("/healthcheck", getHealthcheck) // Healthcheck

	// API
	apiV1 := fiber.New()
	showmasterAPP.Mount("/api/v1", apiV1)
	// Websocket
	apiV1.Get("/ws", websocket.New(a.WebsocketConnection))
	// Login
	apiV1.Post("/login", a.login)                      // <- Email&Password&DeviceID || -> returns new session token
	apiV1.Delete("/logout", a.logout)                  // <- Token, deletes token
	apiV1.Post("/checkLogin", a.checkIfUserIsLoggedIn) // <- Token&DeviceID || -> Bool&Perms, checks if the session is valid and returns the users permissions
	// Admin
	// Web
	showmasterAPP.Static("/", "./frontend/dist/")

	// Start Web
	mst.ElapsedTime()
	log.Println("Started Showmaster V3")
	err = showmasterAPP.Listen(addr)
	if err != nil {
		log.Fatalf("Error starting showmaster v3: %v\n", err)
	}
}
