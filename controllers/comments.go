package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type Comment struct {
	gorm.Model
	Title   string
	Content string

	User      User
	UserID    int
	Process   Process
	ProcessID int

	UnixCreatedAt int64
	UnixUpdatedAt int64
	UnixDeletedAt int64
}

func NewComment(c *fiber.Ctx) error {
	var comment Comment
	db := database.DBConn

	err := c.BodyParser(&comment)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "json_decode_error",
		})
	}

	comment.UnixCreatedAt = time.Now().Unix()
	comment.UserID = RetrieveUserID(c)

	fmt.Printf("[NewComment::comment]: %v", comment)

	db.Create(&comment)
	return c.JSON(comment)
}

// Function that updates all fields of a comment but keeps the current `UnixCreatedAt`
func UpdateComment(c *fiber.Ctx) error {
	var comment Comment
	db := database.DBConn

	err := c.BodyParser(&comment)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "json_decode_error",
		})
	}

	comment.UserID = RetrieveUserID(c)
	comment.UnixUpdatedAt = time.Now().Unix()

	// The following lines return a 403 Forbidden error if the user is not the owner of the comment
	var oldComment Comment
	db.First(&oldComment, comment.ID)
	if oldComment.UserID != comment.UserID {
		return c.Status(403).JSON(&fiber.Map{
			"success": false,
			"cause":   "not_owner",
		})
	}

	fmt.Printf("[UpdateComment::comment]: %v", comment)

	// db.Model(&comment).Omit("UnixCreatedAt").Updates(comment)
	db.Model(&comment).Updates(comment)
	return c.JSON(comment)
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
