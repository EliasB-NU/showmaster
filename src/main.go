package main

import (
	"log"
	"showmaster/src/config"
	"showmaster/src/database"
	"showmaster/src/util"
	"showmaster/src/web"
)

func main() {
	// Start timer
	var mst util.MST
	mst.StartTimer()

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Print("Starting showmaster...")

	// Checks

	// Config
	var CFG = config.GetConfig()

	// Database
	var PSQL = database.GetPSQLDatabase(CFG)
	err := database.InitPSQLDatabase(PSQL)
	if err != nil {
		log.Fatal("Error initializing PSQL database: ", err)
	}

	// Takes the longest, so time is measured here
	var REDIS = database.GetRedisDatabase(CFG)

	// Routines
	util.DeleteOldSessions(PSQL)
	util.DeleteSoftDeletedUserKeys(PSQL)

	// Web
	web.InitWeb(CFG, PSQL, REDIS, &mst)
}
