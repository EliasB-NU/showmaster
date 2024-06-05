package web

import (
	"fmt"
	"net/http"
	"showmaster/src/config"
	"showmaster/src/database"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// All important vars

var ( // All cross Project needed vars
	// Websocket stuff
	WS       *websocket.Conn
	clients  = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Web stuff
	refreshMSG      string           = "refresh"             // Some overengeniered msg
	CFG             config.CFG       = *config.GetConfig()   // The Config
	HighlightedRows []HighlightedRow = initHighlightedRows() // To get all tables with their proberties

	server = fmt.Sprintf("%s:%d", CFG.Website.Host, CFG.Website.Port) // Host & Port
	mux    = http.NewServeMux()                                       // Mux server
	c      = cors.New(cors.Options{                                   // Cors Policy
		AllowedOrigins:       []string{"*"},
		AllowedMethods:       []string{"GET", "POST"},
		AllowPrivateNetwork:  true,
		OptionsSuccessStatus: 200,
	})

	// Stopwatch Stuff
	bu bUpdate

	// Random stuff
	err error
)

// Structs

type HighlightedRow struct { // Struct for all sites, internal use
	Row   float32
	Table string
	Watch Stopwatch
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

type bUpdate struct { // Stopwatch state change incoming message
	RUNNING bool `json:"running"`
	RESET   bool `json:"reset"`
}

type getData struct { // Stopwatch duration & state outgoing message
	Duration uint64
	Running  bool
}

type NumberData struct { // Highlightedrow incoming update message
	Number float32 `json:"number"`
}

type TablesData struct { // All tables outgoing message
	Table string `json:"table"`
}

type NewTableData struct { // New table incoming message
	Name string `json:"name"`
}

type insert struct { // New insert into a table, incoming message
	Table string  `json:"table"`
	ID    float32 `json:"id"`
	Name  string  `json:"name"`
	Audio string  `json:"audio"`
	Licht string  `json:"licht"`
	PPTX  string  `json:"pptx"`
	Notes string  `json:"notes"`
}

// Side Functions

func initHighlightedRows() []HighlightedRow {
	var (
		hr     []HighlightedRow
		tables []string = database.GetTables()
	)

	for _, table := range tables {
		h := HighlightedRow{
			Row:   -1,
			Table: table,
			Watch: *NewStopwatch(),
		}
		hr = append(hr, h)
	}
	return hr
}
