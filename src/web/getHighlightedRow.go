package web

import (
	"encoding/json"
	"log"
	"net/http"
)

type NumberData struct {
	Number float32 `json:"number"`
}

func GetHighlightedRow(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body
	var data NumberData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the received number
	receivedNumber := data.Number
	log.Printf("Received number: %f\n", receivedNumber)

	// Respond to the client
	response := map[string]interface{}{
		"message": "Number received successfully",
	}

	allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "0.0.0.0")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
	json.NewEncoder(w).Encode(response)

	HighlightedRowID = float32(data.Number)
}
