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
	tables []string   = database.GetTables()

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

	for _, table := range tables { // For prject site
		// Serve Websites
		var urlSubfix string = fmt.Sprintf("/%s/", table)
		mux.Handle(urlSubfix, http.StripPrefix(urlSubfix, http.FileServer(http.Dir("./public/rowSite"))))
		// Handle WebSocket connections
		mux.HandleFunc(fmt.Sprintf("/ws%s", urlSubfix), web.HandleConnections)

		// Entries
		mux.HandleFunc(fmt.Sprintf("/api/highlightedrow%s", urlSubfix), web.GetHighlightedRow)
		mux.HandleFunc(fmt.Sprintf("/api/data%s", urlSubfix), web.HandleData)
		//mux.HandleFunc(fmt.Sprintf("/api/newinsert%s", urlSubfix), web.NewInsert)

		// Timer
		mux.HandleFunc(fmt.Sprintf("/api/stopwatch-update%s", urlSubfix), web.ButtonUpdate)
		mux.HandleFunc(fmt.Sprintf("/api/stopwatch-status%s", urlSubfix), web.GetCurrentTime)
	}
	defer web.WS.Close()

	// For project select side
	mux.Handle("/", http.FileServer(http.Dir("./public/projectSite")))
	mux.HandleFunc("/api/gettables", web.HandleTables)

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
