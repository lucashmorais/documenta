package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hlandau/passlib"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type User struct {
	gorm.Model
	Name      string `json: "name"`
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
	Title     string `json: "title"`
	Initials  string `json: "initials"`
	Email     string `json: "email"`
	PHash     string `json: "phash"`
	Roles     []Role `gorm:"many2many:user_roles"`
}

type UserSequence struct {
	gorm.Model
	ProcessID int
	Users     []User `gorm:"many2many:user_sequence_users"`
}

func GetUser(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var user User
	db.Preload("Roles").Where("email = ?", "bob@gmail.com").Where("p_hash = ?", "password123").Find(&user)
	// db.Where("Name = ?", "Albert Billford").Find(&user)

	return c.JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var user []User
	db.Preload("Roles").Find(&user)
	// db.Where("Name = ?", "Albert Billford").Find(&user)

	return c.JSON(user)
}

func PostUser(c *fiber.Ctx) error {
	db := database.DBConn
	var oldUser User

	user := struct {
		gorm.Model
		PHash     string `json: "phash"`
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

	// fmt.Printf("[hash from user-provided password]: %s\n", hash)

	user.PHash = hash

	var roleSlice []Role
	db.Find(&roleSlice, user.Roles)
	fmt.Println(roleSlice)

	dbUser := User{PHash: user.PHash, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName, Initials: user.Initials, Roles: roleSlice}
	db.Create(&dbUser)

	return c.JSON(dbUser)
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

func DeleteUsers(c *fiber.Ctx) error {
	db := database.DBConn
	idsToDelete := struct {
		gorm.Model
		Ids []int `json: "ids"`
	}{}

	err := c.BodyParser(&idsToDelete)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	fmt.Println("[idsToDelete]", idsToDelete)
	fmt.Println("[idsToDelete.Ids]", idsToDelete.Ids)
	db.Delete(&User{}, idsToDelete.Ids)

	return c.SendStatus(200)
}

func GetUserSequences(c *fiber.Ctx) error {
	db := database.DBConn
	var userSequences []UserSequence

	driver := db.Preload("Users")
	processID := c.Query("processID")

	// TODO: Check whether we should convert processID to a number before this
	if processID != "" {
		driver = driver.Where("process_id = ?", processID)
	}

	driver.Find(&userSequences)

	return c.JSON(userSequences)
}

func PostUserSequence(c *fiber.Ctx) error {
	db := database.DBConn
	var userSequence UserSequence

	err := c.BodyParser(&userSequence)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	db.Create(&userSequence)

	return c.JSON(userSequence)
}
