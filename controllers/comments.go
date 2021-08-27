package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type Comment struct {
	gorm.Model
	Title   string `json: "title"`
	Content string `json: "content"`

	User      User
	UserID    int `json: userID`
	Process   Process
	ProcessID int `json:processID`

	UnixCreatedAt int64
	UnixUpdatedAt int64
	UnixDeletedAt int64
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

func GetComments(c *fiber.Ctx) error {
	db := database.DBConn
	var comments []Comment
	driver := db.Preload("User").Preload("Process")

	processID := c.Query("processID")
	userID := c.Query("userID")

	if processID != "" {
		i, err := strconv.Atoi(processID)
		if err == nil {
			driver = driver.Where("process_id = ?", i)
		}
	}

	if userID != "" {
		i, err := strconv.Atoi(userID)
		if err == nil {
			driver = driver.Where("user_id = ?", i)
		}
	}

	driver.Find(&comments)
	return c.JSON(comments)
}
