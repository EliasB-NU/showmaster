package database

import (
	"database/sql"
	"fmt"
	"log"
	"showmaster/src/config"

	_ "github.com/lib/pq"
)

var CFG config.CFG = *config.GetConfig()

func InitDB() *sql.DB {
	// Connect to PostgreSQL database
	psql := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", CFG.DB.Username, CFG.DB.Password, CFG.DB.Host, CFG.DB.Port, CFG.DB.Database)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error opening Database: %v\n", err)
	}

	return db
}

func InitalCheckup() {
	psql := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", CFG.DB.Username, CFG.DB.Password, CFG.DB.Host, CFG.DB.Port, CFG.DB.Database)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error opening Database: %v\n", err)
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id FROM %s", CFG.ProjectName)
	_, err = db.Query(sql)
	if err != nil {
		log.Printf("Error with table: %v\n", err)
		createTable(db)
	}
}

func createTable(db *sql.DB) {
	log.Println("Creating table ...")
	execSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
    	id double precision,
    	name text COLLATE pg_catalog."default",
    	audio text COLLATE pg_catalog."default",
    	licht text COLLATE pg_catalog."default",
    	pptx text COLLATE pg_catalog."default",
    	notes text COLLATE pg_catalog."default"
	)`, CFG.ProjectName)
	_, err := db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error creating table: %v\n", err)
	}
	log.Println("Table created, starting programm")
}
