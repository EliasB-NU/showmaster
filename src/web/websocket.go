package web

import (
	"github.com/gofiber/websocket/v2"
	"log"
	"showmaster/src/database"
)

func (a *API) WebsocketConnection(c *websocket.Conn) {
	a.Clients[c] = true

	type DataElement struct {
		ID          uint64  `json:"id"`
		SceneID     float64 `json:"scene_id"`
		SceneName   string  `json:"scene_name"`
		Audio       string  `json:"audio"`
		Light       string  `json:"light"`
		Video       string  `json:"video"`
		Description string  `json:"description"`
	}
	var (
		activeEventID uint64
		data          []DataElement
		scenes        []database.SceneEntrie

		err error
	)

	for event := range a.LoadedEvents {
		if a.LoadedEvents[event] {
			activeEventID = event
			return
		}
		continue
	}

	err = a.DB.Where("event_id = ?", activeEventID).Find(&scenes).Error
	if err != nil {
		log.Printf("Error getting scenes: %v\n", err)
		return
	}
	for _, scene := range scenes {
		var element DataElement
		element.ID = scene.ID
		element.SceneID = scene.SceneID
		element.SceneName = scene.SceneName
		element.Audio = scene.Audio
		element.Light = scene.Light
		element.Video = scene.Video
		element.Description = scene.SceneDescription

		data = append(data, element)
	}

	err = c.WriteJSON(data)
	if err != nil {
		log.Printf("Error writing data: %v\n", err)
		return
	}

	for {
		_, p, err := c.ReadMessage()
		log.Println(string(p))
		if err != nil {
			break
		}
	}
}

func (a *API) SendMessage(msg []byte) {
	for client := range a.Clients {
		if a.Clients[client] != true {
			continue
		}
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			if err := client.Close(); err != nil {
				log.SetFlags(log.LstdFlags | log.Lshortfile)
				log.Printf("Unable to close websocket connection: %d\n", err)
			}
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error sending message: %d\n", err)
			log.Println("Deleting client ...")
			delete(a.Clients, client)
		}
	}
}
