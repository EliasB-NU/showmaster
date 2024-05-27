package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	"showmaster/src/config"
	"showmaster/src/database"
	"showmaster/src/web"
)

var (
	CFG    config.CFG = *config.GetConfig()
	server            = fmt.Sprintf("%s:%d", CFG.Website.Host, CFG.Website.Port)
	mux               = http.NewServeMux()
)

func main() {
	database.InitalCheckup()

	// Serve Websites
  mux.Handle("/", http.FileServer(http.Dir("./public")))
  mux.Handle("/admin", http.FileServer(http.Dir("./admin")))

	// Handle WebSocket connections
	mux.HandleFunc("/ws", web.HandleConnections)
	go web.AutoUpdate()
	defer web.WS.Close()

	// Handle update on highlighted row
	mux.HandleFunc("/api/highlightedrow", web.GetHighlightedRow)

	// For updating data
	mux.HandleFunc("/api/data", web.HandleData)

	// New Insert
  mux.HandleFunc("/api/newinsert", web.NewInsert)

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
