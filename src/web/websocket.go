package web

import (
	"backend/src/config"
	"backend/src/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	clients  = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var CFG config.CFG = *config.GetConfig()

type Message struct {
	Rows struct {
		ID    float32 `json:"id"`
		Name  string  `json:"name"`
		Audio string  `json:"audio"`
		Licht string  `json:"licht"`
		PPTX  string  `json:"pptx"`
		Notes string  `json:"notes"`
	} `json:"row"`
	Highlighted bool `json:"highlighted"`
}

var HighlightedRowID float32 = 6

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register new client
	clients[ws] = true

	// Send all rows
	SendAllRows(ws)
}

// Update the sendAllRows function to include the highlighted row ID
func SendAllRows(ws *websocket.Conn) {
	// Fetch all rows from the database
	var db sql.DB = *database.InitDB()
	defer db.Close()
	rows, err := db.Query("SELECT id, name, audio, licht, pptx, notes FROM test")
	if err != nil {
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

	// Send JSON data to the client
	if err := ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Printf("Error writing JSON data to client: %v", err)
		return
	}
}
