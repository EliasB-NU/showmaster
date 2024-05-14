package main

import (
	"backend/src/config"
	"backend/src/database"
	"backend/src/web"
	"fmt"
	"log"
	"net/http"
)

var (
	CFG config.CFG = *config.GetConfig()
)

func main() {
	database.InitalCheckup()

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Handle WebSocket connections
	http.HandleFunc("/ws", web.HandleConnections)
	go web.AutoUpdate()
	defer web.WS.Close()

	// Handle update on highlighted row
	http.HandleFunc("/api/highlightedrow", web.GetHighlightedRow)

	// For updating data
	http.HandleFunc("/api/data", web.HandleData)

	// Start the server
	server := fmt.Sprintf("%s:%d", CFG.Website.Host, CFG.Website.Port)
	log.Printf("Server is running on %s\n", server)
	err := http.ListenAndServe(server, nil)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error starting server: %d\n", err)
	}
}
