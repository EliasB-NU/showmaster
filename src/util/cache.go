package util

import (
	"database/sql"
	"fmt"
	"log"
)

// USERS list of all users currently registered
var USERS []string

// PROJECTS list of all projects currently registered
var PROJECTS []string

// CacheUsers caches the users into a var stored in this file, will be called after every user action
func CacheUsers(db *sql.DB) {
	var (
		execSQL = fmt.Sprintf(`SELECT name FROM showmaster.users`)

		err error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error querying rows from showmaster.users: %d\n", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Processing the results
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.SetFlags(log.LstdFlags & log.Lshortfile)
			log.Printf("Error scanning row: %d\n", err)
		}
		USERS = append(USERS, name)
	}

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error with rows: %d\n", err)
	}
}

// CacheProjects caches the projects into a var stored in this file, will be called after every user action
func CacheProjects(db *sql.DB) {
	var (
		execSQL = fmt.Sprintf(`SELECT name FROM showmaster.projects`)

		err error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error querying rows from showmaster.users: %d\n", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Processing the results
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.SetFlags(log.LstdFlags & log.Lshortfile)
			log.Printf("Error scanning row: %d\n", err)
		}
		PROJECTS = append(PROJECTS, name)
	}

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error with rows: %d\n", err)
	}
}
