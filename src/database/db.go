package database

import (
	"backend/src/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var CFG config.CFG = *config.GetConfig()

func InitDB() *sql.DB {
	var err error
	// Connect to PostgreSQL database
	psql := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", CFG.DB.Username, CFG.DB.Password, CFG.DB.Host, CFG.DB.Port, CFG.DB.Daatabase)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error opening Database: %d\n", err)
	}

	return db
}
