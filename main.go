package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/go_fiber/book"
	"github.com/lucashmorais/go_fiber/database"
	"github.com/lucashmorais/go_fiber/products"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "books.db")
	if err != nil {
		panic("Failed to connect to the Database")
	}

	fmt.Println("Database connection successfully established")

	database.DBConn.AutoMigrate(&book.Book{})
	database.DBConn.AutoMigrate(&products.Product{})
	fmt.Println("DB auto-migration was set up")
}

func setupRouter(app *fiber.App) {
	app.Get("/api/v1/book", book.GetBooks)
	app.Get("/api/v1/book/:id", book.GetBook)
	app.Post("/api/v1/book/:id", book.NewBook)
	app.Delete("/api/v1/book/:id", book.DeleteBook)
	app.Get("/api/v1/products", products.GetProducts)
	app.Get("/api/v1/products/:id", products.GetProduct)
	app.Post("/api/v1/products/:id", products.NewProduct)
	app.Delete("/api/v1/products/:id", products.DeleteProduct)
	app.Put("/api/v1/products/:id/quantity/increase", products.IncreaseQuantity)
	app.Put("/api/v1/products/:id/quantity/decrease", products.DecreaseQuantity)
	// app.Get("/", helloWorld)
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		// ServerHeader:  "Fiber",
	})

	initDatabase()

	// This is only executed at the end of the program,
	// for closing the DB connections
	defer database.DBConn.Close()

	setupRouter(app)

	app.Static("/", "./public")

	app.Listen(":3123")
}
