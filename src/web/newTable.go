package web

import (
	"encoding/json"
	"net/http"
	"showmaster/src/database"
)

func NewTable(w http.ResponseWriter, req *http.Request) {
	var data NewTableData
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond to the client
	response := map[string]interface{}{
		"message": "Number received successfully",
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
