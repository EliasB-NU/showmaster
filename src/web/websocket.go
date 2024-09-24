package web

import (
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"log"
)

func WebsocketConnection(c *websocket.Conn) {
	if !c.Locals("allowed").(bool) {
		log.Println("WebSocket upgrade not allowed")
		if err := c.Close(); err != nil {
			log.SetFlags(log.LstdFlags & log.Lshortfile)
			log.Printf("Unable to close websocket connection: %d\n", err)
		}
		return
	}

	clients[c] = true

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

func SendMessage(msg string) {
	jsonMSG, _ := json.Marshal(msg)

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, jsonMSG)
		if err != nil {
			if err := client.Close(); err != nil {
				log.SetFlags(log.LstdFlags & log.Lshortfile)
				log.Printf("Unable to close websocket connection: %d\n", err)
			}
			delete(clients, client)
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error sending message: %v\n", err)
			log.Println("Deleting client ...")
		}
	}
}
