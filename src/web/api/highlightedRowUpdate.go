package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func UpdateHighlightedRow(c *fiber.Ctx) error {
	var (
		urlValue   string = utils.CopyString(c.Params("project"))
		project, _        = strings.CutPrefix(urlValue, ":")
		data       HighlightedRowMSG
	)

	if err := c.BodyParser(&data); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error while recieving new HighlightedRow: %v\n", err)
		c.SendString(fmt.Sprintf("Error while recieving new HighlightedRow: %v\n", err))
		return err
	}

	for id := range HighlightedRows {
		if HighlightedRows[id].Table == project {
			HighlightedRows[id].Row = float32(data.Number)
		}
	}

	SendMessage(project + ":refresh")
	return c.SendString("Number revieved succesfully")
}
