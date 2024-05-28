package web

import (
	"encoding/json"
	"log"
	"net/http"
)

type bUpdate struct {
	RUNNING bool `json:"running"`
	RESET   bool `json:"reset"`
}

type getData struct {
	Duration uint64
	Running  bool
}

var (
	err error
	bu  bUpdate

	Watch Stopwatch = *NewStopwatch()
)

func ButtonUpdate(w http.ResponseWriter, req *http.Request) {
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
		Watch.Reset()
		TimerUpdate("reset")
	} else if bu.RUNNING {
		if Watch.ElapsedSeconds() != 0 {
			Watch.Resume()
			TimerUpdate("start")
		} else {
			Watch.Start()
			TimerUpdate("start")
		}
	} else {
		log.Println(Watch.Stop())
		TimerUpdate("stop")
	}
}

func GetCurrentTime(w http.ResponseWriter, req *http.Request) {
	msg := getData{
		Duration: uint64(Watch.ElapsedSeconds()),
		Running:  Watch.running,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error converting time update msg: %v\n", err)
	}

	w.Write(jsonData)
}
