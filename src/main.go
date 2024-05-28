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
	err error

	CFG    config.CFG = *config.GetConfig()
	server            = fmt.Sprintf("%s:%d", CFG.Website.Host, CFG.Website.Port)

	mux = http.NewServeMux()
	c   = cors.New(cors.Options{
		AllowedOrigins:       []string{"*"},
		AllowedMethods:       []string{"GET", "POST"},
		AllowPrivateNetwork:  true,
		OptionsSuccessStatus: 200,
	})
)

func main() {
	database.InitalCheckup()

	// Serve Websites
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	// Handle WebSocket connections
	mux.HandleFunc("/ws", web.HandleConnections)
	go web.AutoUpdate()
	defer web.WS.Close()

	// Entries
	mux.HandleFunc("/api/highlightedrow", web.GetHighlightedRow)
	mux.HandleFunc("/api/data", web.HandleData)
	mux.HandleFunc("/api/newinsert", web.NewInsert)

	// Timer
	mux.HandleFunc("/api/stopwatch-update", web.ButtonUpdate)
	mux.HandleFunc("/api/stopwatch-status", web.GetCurrentTime)

	// Cors
	handler := c.Handler(mux)

	// Start the server
	log.Printf("Server is running on %s\n", server)
	err = http.ListenAndServe(server, handler)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error starting server: %d\n", err)
	}
}
