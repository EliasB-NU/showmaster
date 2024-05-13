package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

var (
	highlightedRow string
	clients        = make(map[*websocket.Conn]bool)
	broadcast      = make(chan string)
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	db *sql.DB
)

type Row struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Audio string `json:"audio"`
	Licht string `json:"licht"`
	PPTX  string `json:"pptx"`
	Notes string `json:"notes"`
}

func init() {
	var err error
	// Connect to PostgreSQL database
	db, err = sql.Open("postgres", "postgres://showmaster:eb12345678sh@192.168.178.35:5432/showmaster?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register new client
	clients[ws] = true

	// Send all rows to the newly connected client
	sendAllRows(ws)

	for {
		// Listen for new messages from the client
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error: %v", err)
			delete(clients, ws)
			return
		}
	}
}

func handleMessages() {
	for {
		// Wait for a new message to be broadcasted
		newRow := <-broadcast
		highlightedRow = newRow

		// Send the new highlighted row to all connected clients
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(newRow))
			if err != nil {
				log.Printf("Error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func sendAllRows(ws *websocket.Conn) {
	rows, err := db.Query("SELECT id, name, audio, licht, pptx, notes FROM test")
	if err != nil {
		log.Printf("Error fetching rows: %v", err)
		return
	}
	defer rows.Close()

	var allRows []Row
	for rows.Next() {
		var row Row
		if err := rows.Scan(&row.ID, &row.Name, &row.Audio, &row.Licht, &row.PPTX, &row.Notes); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		allRows = append(allRows, row)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return
	}

	// Convert the rows to JSON
	jsonData, err := json.Marshal(allRows)
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

func main() {
	defer db.Close()

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Handle WebSocket connections
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server
	log.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
