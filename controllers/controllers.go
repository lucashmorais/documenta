package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type Book struct {
	gorm.Model
	Title  string `json: "title"`
	Author string `json: "author"`
	Rating int    `json: "rating"`
}

type Process struct {
	gorm.Model
	Title     string `json: "title"`
	Author    string `json: "author"`
	Summary   string `json: "summary"`
	Reference int    `json: "reference"`
	Center    Center
	Type      Type
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
	// ID       int
	Name     string
	Initials string
	Segment  Segment
}

type Segment struct {
	Name string
	//Add any number of permissions
}

type Permission struct {
	Name        string
	Description string
}

type Center struct {
	gorm.Model
	ShortName string
	Name      string
}

type Type struct {
	gorm.Model
	ID          int
	Description string
}

// It seems that Go will only automatically export symbols starting with a capital letter
func NewProcess(c *fiber.Ctx) error {
	db := database.DBConn
	var book Book

	// book.Author = "George Orwell"
	// book.Title = "1984"
	// book.Rating = 4

	// This does not ensure that all Book fields are received nor that
	// requests containing extraneous information is discarded
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).SendString("Request did not include information about the new book")
	}

	db.Create(&book)
	return c.JSON(book)

	// return c.SendString("Adds a new book")
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

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book Book
	db.First(&book, id)
	if book.Title == "" {
		return c.Status(500).SendString("Book was not found")
	}
	title := book.Title
	db.Delete(&book)
	return c.SendString("The book " + title + " was successfully deleted")
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
