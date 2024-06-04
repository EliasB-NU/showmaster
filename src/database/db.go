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
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error connecting to database: %v\n", err)
		return
	}
}

func CreateTable(table string) {
	psql := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", CFG.DB.Username, CFG.DB.Password, CFG.DB.Host, CFG.DB.Port, CFG.DB.Database)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error opening Database: %v\n", err)
		panic(err)
	}
	defer db.Close()

	log.Printf("Creating table %s\n", table)
	execSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id double precision,
		name text,
		audio text,
		licht text,
		pptx text,
		notes text
	)`, table)
	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error creating table: %v\n", err)
	}
	log.Println("Table created")
}
