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
	log.SetFlags(log.LstdFlags & log.Lshortfile)
	log.Println("Starting ShowMaster V3 ...")

	// Load config
	CFG = config.GetConfig()

	// Database
	database.InitDB(CFG)     // Migrates the whole DB and does an initial connection check
	DB = database.GetDB(CFG) // Init the DB var

	// Cache
	util.CacheProjects(DB) // Initial cache of projects
	util.CacheUsers(DB)    // Initial cache of users

	// Web
	w := web.WEB{
		DB:  DB,
		CFG: CFG,
	}
	w.InitWeb()

	log.Println("Started ShowMaster V3")
}
