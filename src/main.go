package main

import (
	"fmt"
	"log"
	"net/http"
	"showmaster/src/config"
	"showmaster/src/database"
	"showmaster/src/web"
)

var (
	CFG config.CFG = *config.GetConfig()
)

func main() {
	database.InitalCheckup()

	// Serve Websites
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/admin", http.FileServer(http.Dir("./admin")))

	// Handle WebSocket connections
	http.HandleFunc("/ws", web.HandleConnections)
	go web.AutoUpdate()
	defer web.WS.Close()

	// Handle update on highlighted row
	http.HandleFunc("/api/highlightedrow", web.GetHighlightedRow)

	// For getting data
	http.HandleFunc("/api/data", web.HandleData)

	// New Insert
	http.HandleFunc("/api/newinsert", web.NewInsert)

	// Start the server
	server := fmt.Sprintf("%s:%d", CFG.Website.Host, CFG.Website.Port)
	log.Printf("Server is running on %s\n", server)
	err := http.ListenAndServe(server, nil)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error starting server: %d\n", err)
	}
}
