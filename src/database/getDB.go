package database

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"showmaster/src/config"
)

// GetDB Gets you a connection of the Database
func GetDB(cfg *config.CFG) *sql.DB {
	var (
		dbURI = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s TimeZone=%s",
			cfg.DB.Host,
			cfg.DB.Username,
			cfg.DB.Database,
			cfg.DB.Password,
			cfg.DB.TimeZone)

		err error = nil
	)

	// Open connection to database
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error connecting to database: %d\n", err)
	} else {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Successfully connected to database")
	}

	database, _ := db.DB()

	return database
}

// InitDB Create schema and tables (if not exist)
func InitDB(cfg *config.CFG) {
	var (
		db  = GetDB(cfg)
		err error
	)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Create schema
	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS showmaster;")
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error creating schema: %d\n", err)
	} else {
		log.Println("Successfully created schema")
	}

	// Create projects table
	execSqlProjects := fmt.Sprint(`
		CREATE TABLE IF NOT EXISTS showmaster.projects (
		    id SERIAL PRIMARY KEY,
		    name TEXT,
		    projecttable TEXT,
		    creator TEXT,
		    timer interval
		);`)

	_, err = db.Exec(execSqlProjects)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error creating table: %d\n", err)
	} else {
		log.Println("Successfully created table : projects")
	}

	// Create users table
	execSqlUsers := fmt.Sprint(`
		CREATE TABLE IF NOT EXISTS showmaster.users (
		    id SERIAL PRIMARY KEY,
		    name TEXT,
		    email TEXT,
		    password TEXT,
		    permlvl integer
		);`)

	_, err = db.Exec(execSqlUsers)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error creating table: %d\n", err)
	} else {
		log.Println("Successfully created table : users")
	}
}
