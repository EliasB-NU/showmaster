package web

import (
	"fmt"
	"log"
	"net/http"
	"showmaster/src/database"
)

func InitiateWeb() {
	var tables []string = database.GetTables()
	for _, table := range tables { // For prject site
		// Serve Websites
		var urlSubfix string = fmt.Sprintf("/%s/", table)
		mux.Handle(urlSubfix, http.StripPrefix(urlSubfix, http.FileServer(http.Dir("./public/rowSite")))) // Site with actuall data

		// Handle WebSocket connections
		mux.HandleFunc(fmt.Sprintf("/ws%s", urlSubfix), HandleConnections) // Websocket

		// Entries
		mux.HandleFunc(fmt.Sprintf("/api/highlightedrow%s", urlSubfix), GetHighlightedRow) // New Highlightedrow
		mux.HandleFunc(fmt.Sprintf("/api/data%s", urlSubfix), HandleData)                  // get all rows
		// mux.HandleFunc(fmt.Sprintf("/api/newinsert%s", urlSubfix), web.NewInsert)           // new insert W.I.P.

		// Timer
		mux.HandleFunc(fmt.Sprintf("/api/stopwatch-update%s", urlSubfix), ButtonUpdate)   // Timer button update
		mux.HandleFunc(fmt.Sprintf("/api/stopwatch-status%s", urlSubfix), GetCurrentTime) // Get current time
	}

	// For project select side
	mux.Handle("/", http.FileServer(http.Dir("./public/projectSite"))) // Website
	mux.HandleFunc("/api/gettables", HandleTables)                     // Get all Tables
	mux.HandleFunc("/api/newtable", NewTable)                          // Create new table

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

func InitiateNewSite(tableName string) {
	// Serve Websites
	var urlSubfix string = fmt.Sprintf("/%s/", tableName)
	mux.Handle(urlSubfix, http.StripPrefix(urlSubfix, http.FileServer(http.Dir("./public/rowSite")))) // Site with actuall data

	// Handle WebSocket connections
	mux.HandleFunc(fmt.Sprintf("/ws%s", urlSubfix), HandleConnections) // Websocket

	// Entries
	mux.HandleFunc(fmt.Sprintf("/api/highlightedrow%s", urlSubfix), GetHighlightedRow) // New Highlightedrow
	mux.HandleFunc(fmt.Sprintf("/api/data%s", urlSubfix), HandleData)                  // get all rows
	// mux.HandleFunc(fmt.Sprintf("/api/newinsert%s", urlSubfix), web.NewInsert)           // new insert W.I.P.

	// Timer
	mux.HandleFunc(fmt.Sprintf("/api/stopwatch-update%s", urlSubfix), ButtonUpdate)   // Timer button update
	mux.HandleFunc(fmt.Sprintf("/api/stopwatch-status%s", urlSubfix), GetCurrentTime) // Get current time

}
