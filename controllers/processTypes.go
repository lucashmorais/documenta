package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type ProcessType struct {
	gorm.Model
	Name        string
	Description string
}

func GetProcessTypes(c *fiber.Ctx) error {
	db := database.DBConn
	var types []ProcessType
	db.Find(&types)

	return c.JSON(types)
}
