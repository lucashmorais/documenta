package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type Permission struct {
	gorm.Model
	Summary     string
	Description string
	// Roles       []Role `gorm:"many2many:role_permissions"`
}

func GetPermissions(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var permission []Permission
	db.Find(&permission)
	// db.Where("Name = ?", "Albert Billford").Find(&permission)

	return c.JSON(permission)
}

func PostPermission(c *fiber.Ctx) error {
	db := database.DBConn
	var permission Permission
	err := c.BodyParser(&permission)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	db.Create(&permission)
	return c.JSON(permission)
}
