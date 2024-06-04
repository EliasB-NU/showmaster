package main

import (
	"showmaster/src/database"
	"showmaster/src/web"
)

func main() {
	database.InitalCheckup()

	web.InitiateWeb()
	defer web.WS.Close() // Close the websocket
}
