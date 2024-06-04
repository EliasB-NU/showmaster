package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func ButtonUpdate(w http.ResponseWriter, req *http.Request) {
	var (
		firstPart, _ = strings.CutPrefix(req.RequestURI, "/api/stopwatch-update/")
		rURL, _      = strings.CutSuffix(firstPart, "/")
	)
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err = json.NewDecoder(req.Body).Decode(&bu); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message": "Update recieved successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	if bu.RESET {
		for id := range HighlightedRows {
			if HighlightedRows[id].Table == rURL {
				HighlightedRows[id].Watch.Reset()
			}
		}
		TimerUpdate("reset+"+rURL, rURL)
	} else if bu.RUNNING {
		for id := range HighlightedRows {
			if HighlightedRows[id].Table == rURL {
				if HighlightedRows[id].Watch.ElapsedSeconds() != 0 {
					HighlightedRows[id].Watch.Resume()
					TimerUpdate("start+"+rURL, rURL)
				} else {
					HighlightedRows[id].Watch.Start()
					TimerUpdate("start+"+rURL, rURL)
				}
			}
		}
	} else {
		for id := range HighlightedRows {
			if HighlightedRows[id].Table == rURL {
				HighlightedRows[id].Watch.Stop()
			}
		}
		TimerUpdate("stop+"+rURL, rURL)
	}
}

func GetCurrentTime(w http.ResponseWriter, req *http.Request) {
	var (
		firstPart, _ = strings.CutPrefix(req.RequestURI, "/api/stopwatch-status/")
		rURL, _      = strings.CutSuffix(firstPart, "/")
		msg          getData
	)

	for id := range HighlightedRows {
		if HighlightedRows[id].Table == rURL {
			msg.Duration = uint64(HighlightedRows[id].Watch.ElapsedSeconds())
			msg.Running = HighlightedRows[id].Watch.running
		}
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error converting time update msg: %v\n", err)
	}

	w.Write(jsonData)
}
