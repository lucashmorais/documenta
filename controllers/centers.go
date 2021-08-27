package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type Center struct {
	gorm.Model
	ShortName string
	Name      string
}

// GET Handler for the Center gorm table.
func GetCenters(c *fiber.Ctx) error {
	db := database.DBConn
	var centers []Center
	db.Find(&centers)

	return c.JSON(centers)
}
