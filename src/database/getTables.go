package database

import (
	"database/sql"
	"log"
)

func GetTables() []string {
	var (
		err error
		db  sql.DB = *GetDB()
	)
	defer db.Close()

	query := `
	SELECT table_name 
	FROM information_schema.tables 
	WHERE table_schema = 'showmaster'
	`

	rows, err := db.Query(query)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	// Process the results
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning row: %v", err)
			return nil
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error with rows: %v", err)
	}

	return tables
}
