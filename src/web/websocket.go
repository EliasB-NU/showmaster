package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"showmaster/src/database"
	"strings"

	"github.com/gorilla/websocket"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	var (
		firstPart, _ = strings.CutPrefix(r.RequestURI, "/ws/")
		rURL, _      = strings.CutSuffix(firstPart, "/")
	)
	// Upgrade initial GET request to a WebSocket
	WS, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error initiating websocket connection with client: %v\n", err)
		return
	}

	// Register new client
	clients[WS] = true

	// Send all rows
	SendAllRows(WS, rURL)
}

// Update the sendAllRows function to include the highlighted row ID
func SendAllRows(ws *websocket.Conn, rURL string) {
	// Fetch all rows from the database
	var db sql.DB = *database.InitDB()
	defer db.Close()
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
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error iterating over rows: %v\n", err)
		return
	}

	// Convert the rows and highlighted row ID to JSON
	jsonData, err := json.Marshal(completeMSG)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	// Send JSON data to the client
	if err := ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error writing JSON data to client: %v\n", err)
		return
	}
}

func SendUpdateMessage() {
	jsonData, _ := json.Marshal(refreshMSG)
	// Send the new highlighted row to all connected clients
	for client := range clients {
		err = client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			client.Close()
			delete(clients, client)
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error sending message: %v\n", err)
		}
	}
}

func TimerUpdate(s string, rURL string) {
	var (
		resetMSG, _ = json.Marshal(s)
		startMSG, _ = json.Marshal(s)
		stopMSG, _  = json.Marshal(s)
	)

	for client := range clients {
		if s == "reset+"+rURL {
			err = client.WriteMessage(websocket.TextMessage, resetMSG)
			if err != nil {
				client.Close()
				delete(clients, client)
				log.SetFlags(log.LstdFlags | log.Lshortfile)
				log.Printf("Error sending message: %v\n", err)
			}
		} else if s == "start+"+rURL {
			err = client.WriteMessage(websocket.TextMessage, startMSG)
			if err != nil {
				client.Close()
				delete(clients, client)
				log.SetFlags(log.LstdFlags | log.Lshortfile)
				log.Printf("Error sending message: %v\n", err)
			}
		} else if s == "stop+"+rURL {
			err = client.WriteMessage(websocket.TextMessage, stopMSG)
			if err != nil {
				client.Close()
				delete(clients, client)
				log.SetFlags(log.LstdFlags | log.Lshortfile)
				log.Printf("Error sending message: %v\n", err)
			}
		} else {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Fatalf("Massive error on timer update websocket stuff: %v\n", err)
		}
	}
}
