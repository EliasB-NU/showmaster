package web

import (
	"backend/src/config"
	"backend/src/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	WS       *websocket.Conn
	clients  = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	err error
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

var refreshMSG string = "refresh"

var HighlightedRowID float32 = -1

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	WS, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Register new client
	clients[WS] = true

	// Send all rows
	SendAllRows(WS)
}

// Update the sendAllRows function to include the highlighted row ID
func SendAllRows(ws *websocket.Conn) {
	// Fetch all rows from the database
	var db sql.DB = *database.InitDB()
	defer db.Close()
	sql := fmt.Sprintf("SELECT id, name, audio, licht, pptx, notes FROM %s", CFG.ProjectName)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Error fetching rows: %v\n", err)
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
		log.Printf("Error iterating over rows: %v\n", err)
		return
	}

	// Convert the rows and highlighted row ID to JSON
	jsonData, err := json.Marshal(completeMSG)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	// Send JSON data to the client
	if err := ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Printf("Error writing JSON data to client: %v\n", err)
		return
	}
}

func AutoUpdate() {
	var previousRow float32 = -1
	jsonData, err := json.Marshal(refreshMSG)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		return
	}
	for {
		// Send the new highlighted row to all connected clients
		if previousRow != HighlightedRowID {
			previousRow = HighlightedRowID
			for client := range clients {
				err = client.WriteMessage(websocket.TextMessage, jsonData)
				if err != nil {
					log.SetFlags(log.LstdFlags | log.Lshortfile)
					log.Printf("Error sending message: %v\n", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
