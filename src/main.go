package main

import (
	"backend/src/config"
	"backend/src/database"
	"backend/src/web"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

var (
	CFG    config.CFG = *config.GetConfig()
	server            = fmt.Sprintf("%s:%d", CFG.Website.Host, CFG.Website.Port)
	mux               = http.NewServeMux()
)

func main() {
	database.InitalCheckup()

	// Serve static files
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	// Handle WebSocket connections
	mux.HandleFunc("/ws", web.HandleConnections)
	go web.AutoUpdate()
	defer web.WS.Close()

	// Handle update on highlighted row
	mux.HandleFunc("/api/highlightedrow", web.GetHighlightedRow)

	// For updating data
	mux.HandleFunc("/api/data", web.HandleData)

	// Cors
	handler := cors.Default().Handler(mux)

	// Start the server
	log.Printf("Server is running on %s\n", server)
	err := http.ListenAndServe(server, handler)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error starting server: %d\n", err)
	}
}
