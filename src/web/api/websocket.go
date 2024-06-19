package api

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
)

func WebsocketConnection(c *websocket.Conn) {
	if !c.Locals("allowed").(bool) {
		log.Println("WebSocket upgrade not allowed")
		c.Close()
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
		err = client.WriteMessage(websocket.TextMessage, jsonMSG)
		if err != nil {
			client.Close()
			delete(clients, client)
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error sending message: %v\n", err)
		}
	}
}
