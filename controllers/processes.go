package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/documenta/database"
)

func GetProcesses(c *fiber.Ctx) error {
	db := database.DBConn
	var processes []Process
	db.Find(&processes)

	return c.JSON(processes)
}

func PostProcess(c *fiber.Ctx) error {
	db := database.DBConn
	var process Process
	err := c.BodyParser(&process)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	db.Create(&process)
	return c.JSON(process)
}
