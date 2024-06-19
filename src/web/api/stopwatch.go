package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func StopWatchStatus(c *fiber.Ctx) error {
	var (
		urlValue   string = utils.CopyString(c.Params("project"))
		project, _        = strings.CutPrefix(urlValue, ":")
		data       StopWatchStatusMSG
	)

	for id := range HighlightedRows {
		if HighlightedRows[id].Table == project {
			data.Duration = uint64(HighlightedRows[id].Watch.ElapsedSeconds())
			data.Running = HighlightedRows[id].Watch.Running
		}
	}

	jsonMsg, err := json.Marshal(data)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error converting time update msg: %v\n", err)
		return err
	}

	return c.Send(jsonMsg)
}

func StopWatchUpdate(c *fiber.Ctx) error {
	var (
		urlValue   string = utils.CopyString(c.Params("project"))
		project, _        = strings.CutPrefix(urlValue, ":")
		data       StopWatchUpdateMSG
	)

	if err := c.BodyParser(&data); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error while recieving stopwatch update: %v\n", err)
		c.SendString(fmt.Sprintf("Error while recieving stopwatch update: %v\n", err))
		return err
	}

	switch data.Update {
	case "start":
		for id := range HighlightedRows {
			if HighlightedRows[id].Table == project {
				if HighlightedRows[id].Watch.ElapsedSeconds() != 0 {
					HighlightedRows[id].Watch.Resume()
					SendMessage(project + ":start")
				} else {
					HighlightedRows[id].Watch.Start()
					SendMessage(project + ":start")
				}
			}
		}

	case "stop":
		for id := range HighlightedRows {
			if HighlightedRows[id].Table == project {
				HighlightedRows[id].Watch.Stop()
				SendMessage(project + ":stop")
			}
		}

	case "reset":
		for id := range HighlightedRows {
			if HighlightedRows[id].Table == project {
				HighlightedRows[id].Watch.Reset()
				SendMessage(project + ":reset")
			}
		}

	default:
		return c.SendString("Unclear update status provided!")
	}

	return c.SendString("Update recieved successfully")
}
