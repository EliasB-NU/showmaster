package util

import (
	"database/sql"
	"fmt"
	"log"
)

type ProjectCache struct {
	Table          string
	HighlightedRow float32
	Timer          Stopwatch
}

// USERS list of all users currently registered
var USERS []string

// PROJECTS list of all projects currently registered
var PROJECTS []string

// TABLES list of al tables and their cached data currently registered
var TABLES []ProjectCache

// CacheUsers caches the users into a var stored in this file, will be called after every user action
func CacheUsers(db *sql.DB) {
	var (
		execSQL = fmt.Sprintf(`SELECT name FROM showmaster.users`)

		err error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying rows from showmaster.users: %d\n", err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error iterating over rows: %d\n", err)
	}

	// Processing the results
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning row: %d\n", err)
		}
		USERS = append(USERS, name)
	}
}

// CacheProjects caches the projects into a var stored in this file, will be called after every project action
func CacheProjects(db *sql.DB) {
	var (
		execSQL = fmt.Sprintf(`SELECT name FROM showmaster.projects`)

		err error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying rows from showmaster.users: %d\n", err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error iterating over rows: %d\n", err)
	}

	// Processing the results
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning row: %d\n", err)
		}
		PROJECTS = append(PROJECTS, name)
	}
}

// CacheTables caches the tables into a var stored in this file, will be called on every HighlightedRow/timer update
func CacheTables(db *sql.DB) {
	type d2 struct {
		Table string
		Timer *float64
	}

	var (
		execSQL  = fmt.Sprintf(`SELECT projecttable, timer FROM showmaster.projects`)
		projects []d2
	)
	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying names of all tables: %d\n", err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error iterating over rows: %d\n", err)
	}

	for rows.Next() {
		var n d2
		if err := rows.Scan(&n.Table, &n.Timer); err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning row: %d\n", err)
		}
		projects = append(projects, n)
	}
	for _, p := range projects {
		var d ProjectCache
		d.Table = p.Table
		d.HighlightedRow = 0
		d.Timer = *NewStopwatch()
		if !(p.Timer == nil) {
			d.Timer.SetElapsedSeconds(*p.Timer)
		}

		TABLES = append(TABLES, d)
	}
}

// NewTableCache on the creation of a new project, it caches stored values get added to the cache variable
func NewTableCache(projectTable string) {
	var d ProjectCache
	d.Table = projectTable
	d.HighlightedRow = 0
	d.Timer = *NewStopwatch()

	TABLES = append(TABLES, d)
}
