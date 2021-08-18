package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/documenta/database"
)

func GetProcessTypes(c *fiber.Ctx) error {
	db := database.DBConn
	var types []ProcessType
	db.Find(&types)

	return c.JSON(types)
}
