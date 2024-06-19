package api

import (
	"showmaster/src/database"
	"showmaster/src/util"

	"github.com/gofiber/websocket/v2"
)

var (
	clients = make(map[*websocket.Conn]bool)

	HighlightedRows []HighlightedRow = initHighlightedRows()

	err error
)

type HighlightedRow struct { // Struct for all sites, internal use
	Row   float32
	Table string
	Watch util.Stopwatch
}

type RowMSG struct { // Rows outgoing data
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

type TableMSG struct { // All tables outgoing message
	Table string `json:"table"`
}

type HighlightedRowMSG struct { // Highlightedrow incoming update message
	Number float32 `json:"number"`
}

type StopWatchStatusMSG struct { // Stopwatch duration & state outgoing message
	Duration uint64
	Running  bool
}

type StopWatchUpdateMSG struct { // Stopwatch state change incoming message
	Update string `json:"update"`
}

type NewTableMSG struct { // New table incoming message
	Name string `json:"name"`
}

type NewInsertMSG struct { // New insert into a table, incoming message
	ID    float32 `json:"id"`
	Name  string  `json:"name"`
	Audio string  `json:"audio"`
	Licht string  `json:"licht"`
	PPTX  string  `json:"pptx"`
	Notes string  `json:"notes"`
}

func initHighlightedRows() []HighlightedRow {
	var (
		hr     []HighlightedRow
		tables []string = database.GetTables()
	)

	for _, table := range tables {
		h := HighlightedRow{
			Row:   -1,
			Table: table,
			Watch: *util.NewStopwatch(),
		}
		hr = append(hr, h)
	}
	return hr
}
