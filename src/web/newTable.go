package web

import (
	"encoding/json"
	"net/http"
	"showmaster/src/database"
)

func NewTable(w http.ResponseWriter, req *http.Request) {
	var (
		data           NewTableData
		tables         []string = database.GetTables()
		alreadyInTable bool     = false
	)
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, t := range tables {
		if data.Name == t {
			alreadyInTable = true
		} else {
			continue
		}
	}

	if alreadyInTable {
		// Respond to the client
		response := map[string]interface{}{
			"message": "Already in table",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		// Respond to the client
		response := map[string]interface{}{
			"message": "Table received successfully",
		}
		json.NewEncoder(w).Encode(response)

		database.CreateTable(data.Name)

		InitiateNewSite(data.Name)

		hr := HighlightedRow{
			Row:   -1,
			Table: data.Name,
			Watch: *NewStopwatch(),
		}

		HighlightedRows = append(HighlightedRows, hr)
	}
}
