package main

import (
	"database/sql"
	"log"
	"showmaster/src/config"
	"showmaster/src/database"
	"showmaster/src/util"
	"showmaster/src/web"
)

var DB *sql.DB
var CFG *config.CFG

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting ShowMaster V3 ...")

	// Load config
	CFG = config.GetConfig()

	// Database
	database.InitDB(CFG)     // Migrates the whole DB and does an initial connection check
	DB = database.GetDB(CFG) // Init the DB var
	defer DB.Close()         // Close the DB on shutdown

	// Cache
	util.CacheProjects(DB)                      // Initial cache of projects
	util.CacheUsers(DB)                         // Initial cache of users
	util.CacheTables(DB)                        // Initial Cache of all projects
	err := util.CreateInitialAdminUser(DB, CFG) // Creates the initial admin user from the docker compose file values
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error creating initial admin user: %v", err)
	}

	// Web
	web.InitWeb(DB, CFG)
}
