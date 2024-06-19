package main

import (
	"showmaster/src/database"
	"showmaster/src/web"
)

func main() {
	database.InitalCheckup()

	web.SetLogLevel()           // Log level
	go web.CheckForNewRowSite() // because of import cycle problem, periodic check if there is a new table
	web.StartTheWeb()           // Inits all the api endpoints and sites
}
