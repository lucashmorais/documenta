package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hlandau/passlib"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

func GetUser(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var user User
	db.Where("email = ?", "bob@gmail.com").Where("p_hash = ?", "password123").Find(&user)
	// db.Where("Name = ?", "Albert Billford").Find(&user)

	return c.JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var user []User
	db.Find(&user)
	// db.Where("Name = ?", "Albert Billford").Find(&user)

	return c.JSON(user)
}

func PostUser(c *fiber.Ctx) error {
	db := database.DBConn
	var user, oldUser User
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	db.Where("email = ?", user.Email).Find(&oldUser)

	if oldUser.Email != "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "email_was_taken",
		})
	}

	hash, err := passlib.Hash(user.PHash)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "cannot_hash_password",
		})
	}

	fmt.Printf("[hash from user-provided password]: %s", hash)

	user.PHash = hash

	db.Create(&user)
	return c.JSON(user)
}

func PutUser(c *fiber.Ctx) error {
	db := database.DBConn
	user := struct {
		gorm.Model
		Email     string `json: "email"`
		FirstName string `json: "firstName"`
		LastName  string `json: "lastName"`
		Initials  string `json: "initials"`
		Roles     []int  `json: "roles"`
	}{}

	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	dbUser := User{}

	var roleSlice []Role

	db.Where("id = ?", user.ID).First(&dbUser)
	db.Find(&roleSlice, user.Roles)

	db.Model(&dbUser).Update(user)
	db.Model(&dbUser).Association("Roles").Replace(roleSlice)
	return c.JSON(user)
}
