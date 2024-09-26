package web

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
)

func (a *API) getUsers(c *fiber.Ctx) error {
	users, err := database.GetUsers(a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error getting users: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

type updateUser struct {
	Name    string `json:"name"`
	Permlvl int    `json:"permlvl"`
}

func (a *API) updateUser(c *fiber.Ctx) error {
	var (
		data updateUser
		err  error
	)

	if err = c.BodyParser(&data); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing body: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = database.UpdateUser(data.Name, data.Permlvl, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error updating user: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	util.CacheUsers(a.DB)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

type deleteStruct struct {
	Name string `json:"name"`
}

func (a *API) deleteUser(c *fiber.Ctx) error {
	var (
		data deleteStruct
		err  error
	)

	if err := c.BodyParser(&data); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing json: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = database.DeleteUser(data.Name, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error deleting user: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	util.CacheUsers(a.DB)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (a *API) deleteProject(c *fiber.Ctx) error {
	var (
		data deleteStruct
		err  error
	)

	if err := c.BodyParser(&data); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error parsing json: %d\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	err = database.DeleteProject(data.Name, a.DB)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error deleting project: %d\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	util.CacheProjects(a.DB)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
