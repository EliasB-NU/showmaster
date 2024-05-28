package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"showmaster/src/config"
	"showmaster/src/database"

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

	refreshMSG       string     = "refresh"
	CFG              config.CFG = *config.GetConfig()
	HighlightedRowID float32    = -1
	Update           bool       = false
)

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

func HandleConnections(w http.ResponseWriter, r *http.Request) {
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
		if previousRow != HighlightedRowID || Update {
			previousRow = HighlightedRowID
			Update = false
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
	}
}

func TimerUpdate(s string) {
	resetMSG, _ := json.Marshal("reset")
	startMSG, _ := json.Marshal("start")
	stopMSG, _ := json.Marshal("stop")

	for client := range clients {
		if s == "reset" {
			err = client.WriteMessage(websocket.TextMessage, resetMSG)
			if err != nil {
				client.Close()
				delete(clients, client)
				log.SetFlags(log.LstdFlags | log.Lshortfile)
				log.Printf("Error sending message: %v\n", err)
			}
		} else if s == "start" {
			err = client.WriteMessage(websocket.TextMessage, startMSG)
			if err != nil {
				client.Close()
				delete(clients, client)
				log.SetFlags(log.LstdFlags | log.Lshortfile)
				log.Printf("Error sending message: %v\n", err)
			}
		} else if s == "stop" {
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
