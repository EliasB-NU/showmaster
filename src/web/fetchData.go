package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"showmaster/src/database"
)

func HandleData(w http.ResponseWriter, r *http.Request) {
	// Fetch all rows from the database
	var db sql.DB = *database.InitDB()
	defer db.Close()
	sql := fmt.Sprintf("SELECT id, name, audio, licht, pptx, notes FROM %s", CFG.ProjectName)
	rows, err := db.Query(sql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error fetching rows: %v", err)
		return
	}
	defer rows.Close()

	// Iterate over the rows and construct Row objects
	var completeMSG []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Rows.ID, &msg.Rows.Name, &msg.Rows.Audio, &msg.Rows.Licht, &msg.Rows.PPTX, &msg.Rows.Notes); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}
		if msg.Rows.ID == HighlightedRowID {
			msg.Highlighted = true
		}
		completeMSG = append(completeMSG, msg)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return
	}

	// Convert the rows and highlighted row ID to JSON
	jsonData, err := json.Marshal(completeMSG)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	// Set response headers
	allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

	// Write JSON response
	w.Write(jsonData)
}
