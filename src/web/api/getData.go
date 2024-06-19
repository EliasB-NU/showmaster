package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"showmaster/src/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func GetData(c *fiber.Ctx) error {
	var (
		urlValue   string  = utils.CopyString(c.Params("project"))
		project, _         = strings.CutPrefix(urlValue, ":")
		db         *sql.DB = database.GetDB()
	)
	defer db.Close()

	sql := fmt.Sprintf("SELECT id, name, audio, licht, pptx, notes FROM %s", project)
	rows, err := db.Query(sql)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error fetching rows: %v\n", err)
		return err
	}
	defer rows.Close()

	// Iterate over the rows and construct Row objects
	var completeMSG []RowMSG
	for rows.Next() {
		var msg RowMSG
		if err := rows.Scan(&msg.Rows.ID, &msg.Rows.Name, &msg.Rows.Audio, &msg.Rows.Licht, &msg.Rows.PPTX, &msg.Rows.Notes); err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning row: %v\n", err)
			continue
		}
		for id := range HighlightedRows {
			if HighlightedRows[id].Table == project {
				if msg.Rows.ID == HighlightedRows[id].Row {
					msg.Highlighted = true
				}
			}
		}
		completeMSG = append(completeMSG, msg)
	}
	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error iterating over rows: %v", err)
		return err
	}

	// Convert the rows and highlighted row ID to JSON
	jsonMSG, err := json.Marshal(completeMSG)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error marshalling JSON: %v\n", err)
		return err
	}

	return c.Send(jsonMSG)
}

func GetTables(c *fiber.Ctx) error {
	var (
		tables      []string = database.GetTables()
		completeMSG []TableMSG
	)

	for i := 0; i < len(tables); i++ {
		zwischenAaal := TableMSG{
			Table: tables[i],
		}
		completeMSG = append(completeMSG, zwischenAaal)
	}
	jsonMSG, err := json.Marshal(completeMSG)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error marshalling JSON: %v\n", err)
		return err
	}

	return c.Send(jsonMSG)
}
