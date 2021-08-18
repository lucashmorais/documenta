package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type Process struct {
	gorm.Model
	Title     string `json: "title"`
	Summary   string `json: "summary"`
	Reference int    `json: "reference"`
	Author    User
	Status    ProcessStatus
	Center    Center
	Type      ProcessType
}

type Comment struct {
	gorm.Model
	Title         string `json: "title"`
	Content       string `json: "content"`
	UserID        int    `json: userID`
	User          User
	ProcessID     int `json:processID`
	Process       Process
	UnixCreatedAt int64
	UnixUpdatedAt int64
	UnixDeletedAt int64
}

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

type Role struct {
	gorm.Model
	Name        string
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	gorm.Model
	Summary     string
	Description string
	// Roles       []Role `gorm:"many2many:role_permissions"`
}

type Center struct {
	gorm.Model
	ShortName string
	Name      string
}

type ProcessStatus struct {
	gorm.Model
	Name        string
	Description string
}

type ProcessType struct {
	gorm.Model
	Name        string
	Description string
}

func NewComment(c *fiber.Ctx) error {
	var comment Comment
	db := database.DBConn

	defaultUser := User{Name: "Albert Billford", Initials: "AB"}
	defaultProcess := Process{Title: "Initial Process"}

	if err := c.BodyParser(&comment); err != nil {
		return c.Status(400).SendString("Request did not include information about the new comment")
	}

	comment.UnixCreatedAt = time.Now().Unix()

	db.Create(&defaultUser)
	comment.User = defaultUser

	db.Create(&defaultProcess)
	comment.Process = defaultProcess

	db.Create(&comment)
	return c.JSON(comment)
}

func UpdateComment(c *fiber.Ctx) error {
	return nil
}

func DeleteComment(c *fiber.Ctx) error {
	return nil
}

func GetComment(c *fiber.Ctx) error {
	return nil
}

func GetCommentsByProcessID(c *fiber.Ctx) error {
	process_id := c.Params("id")

	db := database.DBConn
	var comments []Comment

	//CHECK FIELD NAMES ON DB BROWSER BEFORE QUERYING DATA!!!
	db.Preload("User").Preload("Process").Where("process_id = ?", process_id).Find(&comments)

	return c.JSON(comments)
}

func GetComments(c *fiber.Ctx) error {
	db := database.DBConn
	var comments []Comment
	db.Preload("User").Preload("Process").Where("process_id = ?", 1).Find(&comments)
	db.Preload("User").Preload("Process").Find(&comments)

	return c.JSON(comments)
}
