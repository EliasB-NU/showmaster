package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func GetHighlightedRow(w http.ResponseWriter, r *http.Request) {
	var (
		firstPart, _ = strings.CutPrefix(r.RequestURI, "/api/highlightedrow/")
		rURL, _      = strings.CutSuffix(firstPart, "/")
	)
	// Decode the JSON request body
	var data NumberData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the received number
	log.Printf("Received number from %s: %f\n", rURL, data.Number)

	// Respond to the client
	response := map[string]interface{}{
		"message": "Number received successfully",
	}
	json.NewEncoder(w).Encode(response)

	for id := range HighlightedRows {
		if HighlightedRows[id].Table == rURL {
			HighlightedRows[id].Row = float32(data.Number)
		}
	}
	SendUpdateMessage()
}
