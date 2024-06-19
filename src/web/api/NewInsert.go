package api

import (
	"database/sql"
	"fmt"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func NewTable(c *fiber.Ctx) error {
	var (
		data           NewTableMSG
		tables         []string = database.GetTables()
		alreadyInTable bool     = false
	)

	if err := c.BodyParser(&data); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error while revieving new table: %v\n", err)
		return c.SendString(fmt.Sprintf("Error while revieving new table: %v\n", err))
	}

	for _, t := range tables {
		if data.Name == t {
			alreadyInTable = true
		} else {
			continue
		}
	}

	if alreadyInTable {
		// Respond to the client
		return c.SendString("Table already exists")
	} else {

		database.CreateTable(data.Name)

		hr := HighlightedRow{
			Row:   -1,
			Table: data.Name,
			Watch: *util.NewStopwatch(),
		}

		HighlightedRows = append(HighlightedRows, hr)

		return c.SendString("New table recieved successfully")
	}
}

func NewInsert(c *fiber.Ctx) error {
	var (
		urlValue   string = utils.CopyString(c.Params("project"))
		project, _        = strings.CutPrefix(urlValue, ":")
		i          NewInsertMSG
		db         *sql.DB = database.GetDB()
	)

	if err := c.BodyParser(&i); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error with recieved insert: %v\n", err)
		return c.SendString(fmt.Sprintf("Error with recieved insert: %v\n", err))
	}

	sql := fmt.Sprintf("INSERT INTO %s (id, name, audio, licht, pptx, notes) VALUES (%f, '%s', '%s', '%s', '%s', '%s');", project, i.ID, i.Name, i.Audio, i.Licht, i.PPTX, i.Notes)
	_, err := db.Exec(sql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error inserting new entry: %v\n", err)
	}

	return nil
}
