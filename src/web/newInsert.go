package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"showmaster/src/database"
)

func NewInsert(w http.ResponseWriter, req *http.Request) {
	var (
		i  insert
		db *sql.DB = database.InitDB()
	)
	defer db.Close()

	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&i); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond to the client
	response := map[string]interface{}{
		"message": "Newinsert received successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Insert new entry into database, and tell the client to refresh
	sql := fmt.Sprintf("INSERT INTO %s (id, name, audio, licht, pptx, notes) VALUES (%f, '%s', '%s', '%s', '%s', '%s');", i.Table, i.ID, i.Name, i.Audio, i.Licht, i.PPTX, i.Notes)
	_, err := db.Exec(sql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error inserting new entry: %v\n", err)
	}
}
