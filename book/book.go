package book

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/go_fiber/database"
)

type Book struct {
	gorm.Model
	Title  string `json: "title"`
	Author string `json: "author"`
	Rating int    `json: "rating"`
}

func GetBooks(c *fiber.Ctx) error {
	db := database.DBConn
	var books []Book
	db.Find(&books)
	return c.JSON(books)
	// return c.SendString("All Books")
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book Book
	db.Find(&book, id)
	return c.JSON(book)
	// return c.SendString("Single book")
}

func NewBook(c *fiber.Ctx) error {
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
