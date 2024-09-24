package main

import (
	"database/sql"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
)

var DB *sql.DB

func main() {
	log.SetFlags(log.LstdFlags & log.Lshortfile)
	log.Println("Starting ShowMaster V3 ...")

	// Database
	database.InitDB()     // Migrates the whole DB and does an initial connection check
	DB = database.GetDB() // Init the DB var

	// Cache
	util.CacheProjects(DB) // Initial cache of projects
	util.CacheUsers(DB)    // Initial cache of users

	log.Println("Started ShowMaster V3")
}
