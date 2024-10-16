package web

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
)

func (a *API) getProjects(c *fiber.Ctx) error {
	data, err := database.GetProjects(a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error getting projects: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

func (a *API) newProject(c *fiber.Ctx) error {
	var (
		data database.Project
		err  error
	)

	err = c.BodyParser(&data)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing request: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = database.NewProject(data.Name, data.Creator, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error creating project: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	util.CacheProjects(a.DB)
	util.NewTableCache(data.Name + "table")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (a *API) updateProject(c *fiber.Ctx) error {
	type dataStruct struct {
		OldName string `json:"oldname"`
		NewName string `json:"newname"`
	}
	var (
		data dataStruct
		err  error
	)

	err = c.BodyParser(&data)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing request: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = database.UpdateProject(data.OldName, data.NewName, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error updating project: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	util.CacheProjects(a.DB)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
