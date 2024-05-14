package web

import (
	"backend/src/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type NumberData struct {
	Number float32 `json:"number"`
}

func GetHighlightedRow(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body
	var data NumberData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the received number
	receivedNumber := data.Number
	log.Printf("Received number: %f\n", receivedNumber)

	// Respond to the client
	response := map[string]interface{}{
		"message": "Number received successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	HighlightedRowID = float32(data.Number)

	requestURL := fmt.Sprintf("http://localhost:%d/api/refresh", CFG.Website.Port)
	_, err := http.Get(requestURL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error sending refresh request to clients: %v\n", err)
	}
}

func HandleData(w http.ResponseWriter, r *http.Request) {
	// Fetch all rows from the database
	var db sql.DB = *database.InitDB()
	defer db.Close()
	rows, err := db.Query("SELECT id, name, audio, licht, pptx, notes FROM test")
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
		if err := rows.Scan(&msg.Rows.ID, &msg.Rows.Name, &msg.Rows.Audio, &msg.Rows.Licht, &msg.Rows.PPTX, &msg.Rows.Name); err != nil {
			log.Printf("Error scanning row: %v", err)
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
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	w.Write(jsonData)
}
