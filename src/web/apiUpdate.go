package web

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
	"strings"
)

// Errors
var (
	projectNotFound = errors.New("project not found")
)

func (a *API) getHighlightedRow(c *fiber.Ctx) error {
	var (
		urlValue   = utils.CopyString(c.Params("project"))
		project, _ = strings.CutPrefix(urlValue, ":")

		err error
	)
	if project == "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request without project link ...")
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	p, err := findProject(project)
	if errors.Is(err, projectNotFound) {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request with invalid project link ...")
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	return c.Status(fiber.StatusOK).JSON(p.HighlightedRow)
}

func (a *API) updateHighlightedRow(c *fiber.Ctx) error {
	type returnMSG struct {
		Table          string  `json:"table"`
		HighlightedRow float32 `json:"highlighted_row"`
	}
	var (
		urlValue   = utils.CopyString(c.Params("project"))
		project, _ = strings.CutPrefix(urlValue, ":")

		data float32
		msg  returnMSG
		err  error
	)
	if project == "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request without project link ...")
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	p, err := findProject(project)
	if errors.Is(err, projectNotFound) {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request with invalid project link ...")
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	err = c.BodyParser(&data)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsting body: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	p.HighlightedRow = data

	msg.Table = p.Table
	msg.HighlightedRow = data
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error marshalling data: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}
	SendMessage(bytes)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (a *API) getTimer(c *fiber.Ctx) error {
	type msg struct {
		Status bool    `json:"status"`
		Timer  float64 `json:"timer"`
	}
	var (
		urlValue   = utils.CopyString(c.Params("project"))
		project, _ = strings.CutPrefix(urlValue, ":")

		data msg
		err  error
	)
	if project == "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request without project link ...")
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	p, err := findProject(project)
	if errors.Is(err, projectNotFound) {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request with invalid project link ...")
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	data.Status = p.Timer.Running
	data.Timer = p.Timer.ElapsedSeconds()

	return c.Status(fiber.StatusOK).JSON(data)
}

func (a *API) updateTimer(c *fiber.Ctx) error {
	type returnMSG struct {
		Table       string `json:"table"`
		TimerStatus string `json:"timer_status"`
	}
	var (
		urlValue   = utils.CopyString(c.Params("project"))
		project, _ = strings.CutPrefix(urlValue, ":")

		data string
		msg  returnMSG
		err  error
	)
	if project == "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request without project link ...")
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	p, err := findProject(project)
	if errors.Is(err, projectNotFound) {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request with invalid project link ...")
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	err = c.BodyParser(&data)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing body: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON("")
	}

	switch data {
	case "start":
		if p.Timer.ElapsedSeconds() != 0 {
			p.Timer.Resume()
			msg.TimerStatus = "start"
		} else {
			p.Timer.Start()
			msg.TimerStatus = "start"
		}
	case "stop":
		p.Timer.Stop()
		msg.TimerStatus = "stop"
		err := database.UpdateTimer(p.Timer.ElapsedSeconds(), p.Table, a.DB)
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error updating timer: %d\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON("")
		}
	case "reset":
		p.Timer.Reset()
		err := database.UpdateTimer(0, p.Table, a.DB)
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error updating timer: %d\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON("")
		}
		msg.TimerStatus = "reset"
	default:
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request with invalid status code")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	msg.Table = p.Table

	bytes, err := json.Marshal(&msg)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error marshalling data: %s\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}
	SendMessage(bytes)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func findProject(project string) (util.ProjectCache, error) {
	var dummyData util.ProjectCache
	for _, v := range util.TABLES {
		if v.Table == project {
			return v, nil
		}
	}
	return dummyData, projectNotFound
}
