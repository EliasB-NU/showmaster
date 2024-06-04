package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"showmaster/src/database"
	"strings"
)

func HandleData(w http.ResponseWriter, r *http.Request) {
	var (
		firstPart, _        = strings.CutPrefix(r.RequestURI, "/api/data/")
		rURL, _             = strings.CutSuffix(firstPart, "/")
		db           sql.DB = *database.InitDB()
	)
	defer db.Close()

	// Fetch all rows from the database
	sql := fmt.Sprintf("SELECT id, name, audio, licht, pptx, notes FROM %s", rURL)
	rows, err := db.Query(sql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error fetching rows: %v\n", err)
		return
	}
	defer rows.Close()

	// Iterate over the rows and construct Row objects
	var completeMSG []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Rows.ID, &msg.Rows.Name, &msg.Rows.Audio, &msg.Rows.Licht, &msg.Rows.PPTX, &msg.Rows.Notes); err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning row: %v\n", err)
			continue
		}
		for id := range HighlightedRows {
			if HighlightedRows[id].Table == rURL {
				if msg.Rows.ID == HighlightedRows[id].Row {
					msg.Highlighted = true
				}
			}
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
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	// Write JSON response
	w.Write(jsonData)
}

func HandleTables(w http.ResponseWriter, req *http.Request) {
	var (
		tables      []string = database.GetTables()
		completeMSG []TablesData
	)
	for i := 0; i < len(tables); i++ {
		zwischenAaal := TablesData{
			Table: tables[i],
		}
		completeMSG = append(completeMSG, zwischenAaal)
	}
	jsonData, err := json.Marshal(completeMSG)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	w.Write(jsonData)
}
