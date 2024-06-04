package web

import (
	"net/http"
	"showmaster/src/config"
	"showmaster/src/database"

	"github.com/gorilla/websocket"
)

type bUpdate struct { // Stopwatch state change incoming message
	RUNNING bool `json:"running"`
	RESET   bool `json:"reset"`
}

type getData struct { // Stopwatch duration & state outgoing message
	Duration uint64
	Running  bool
}

type HighlightedRow struct { // Struct for all sites, internal use
	Row   float32
	Table string
	Watch Stopwatch
}

type NumberData struct { // Highlightedrow incoming update message
	Number float32 `json:"number"`
}

type Message struct { // Rows outgoing data
	Rows struct {
		ID    float32 `json:"id"`
		Name  *string `json:"name"`
		Audio *string `json:"audio"`
		Licht *string `json:"licht"`
		PPTX  *string `json:"pptx"`
		Notes *string `json:"notes"`
	} `json:"row"`
	Highlighted bool `json:"highlighted"`
}

var (
	// Websocket stuff
	WS       *websocket.Conn 
	clients  = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Web stuff
	refreshMSG      string           = "refresh"
	Tables          []string         = database.GetTables()
	CFG             config.CFG       = *config.GetConfig()
	HighlightedRows []HighlightedRow = initHighlightedRows()

	// Stopwatch Stuff
	bu bUpdate

	// Random stuff
	err error
)
