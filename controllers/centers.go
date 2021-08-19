package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/documenta/database"
)

func GetCenters(c *fiber.Ctx) error {
	db := database.DBConn
	var centers []Center
	db.Find(&centers)

	return c.JSON(centers)
}
