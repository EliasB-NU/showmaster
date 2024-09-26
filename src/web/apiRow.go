package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"showmaster/src/database"
	"strings"
)

func (a *API) getRows(c *fiber.Ctx) error {
	var (
		urlValue   = utils.CopyString(c.Params("project"))
		project, _ = strings.CutPrefix(urlValue, ":")
	)
	if project == "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request without project link ...")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	data, err := database.GetRows(project, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error getting rows: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

func (a *API) newRow(c *fiber.Ctx) error {
	var (
		urlValue   = utils.CopyString(c.Params("project"))
		project, _ = strings.CutPrefix(urlValue, ":")
		data       database.Row
		err        error
	)
	if project == "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request without project link ...")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = c.BodyParser(&data)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing data: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = database.NewRow(project, data, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error creating row: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (a *API) updateRow(c *fiber.Ctx) error {
	var (
		urlValue   = utils.CopyString(c.Params("project"))
		project, _ = strings.CutPrefix(urlValue, ":")
		data       database.Row
		err        error
	)
	if project == "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request without project link ...")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = c.BodyParser(&data)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing data: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = database.UpdateRow(project, data, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error creating row: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (a *API) deleteRow(c *fiber.Ctx) error {
	var (
		urlValue   = utils.CopyString(c.Params("project"))
		project, _ = strings.CutPrefix(urlValue, ":")
		data       int
		err        error
	)
	if project == "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Request without project link ...")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = c.BodyParser(&data)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing data: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = database.DeleteRow(project, data, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error deleting row: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
