package main

import (
	"context"
	"log"
	"showmaster/backend/config"
	"showmaster/backend/database"
	"showmaster/backend/util"
	"showmaster/backend/web"
)

func main() {
	// Create a new stopwatch to measure startup time
	var mst util.MST
	mst.StartTimer()

	// Showmaster V3
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Showmaster V3 ...")

	// Create systemwide context
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Get Config
	var cfg = config.GetConfig()

	// Get Databases
	var (
		psql = database.GetPSQL(cfg)
		rdb  = database.GetRedis(cfg)
	)
	// Init config
	err := database.InitPSQLDatabase(psql)
	if err != nil {
		log.Fatalf("InitPSQLDatabase failed: %v\n", err)
	}

	// Routines
	util.DeleteOldSessions(psql)
	util.DeleteSoftDeletedUserKeys(psql)

	// Init Web
	web.InitWeb(psql, rdb, ctx, cfg, &mst)
}
