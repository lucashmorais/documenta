package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

func GetRoles(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var role []Role
	db.Find(&role)
	// db.Where("Name = ?", "Albert Billford").Find(&role)

	return c.JSON(role)
}

func PostRole(c *fiber.Ctx) error {
	db := database.DBConn
	role := struct {
		gorm.Model
		Name        string `json: "name"`
		Description string `json: "description"`
		Permissions []int  `json: "permissions"`
	}{}

	err := c.BodyParser(&role)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	fmt.Println(role)

	var permissionSlice []Permission
	db.Find(&permissionSlice, role.Permissions)
	fmt.Println(permissionSlice)

	dbRole := Role{Name: role.Name, Description: role.Description, Permissions: permissionSlice}

	db.Create(&dbRole)
	return c.JSON(role)
}