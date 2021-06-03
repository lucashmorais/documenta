package main

import (
	"fmt"
	"time"

	"github.com/lucashmorais/documenta/controllers"
	"github.com/lucashmorais/documenta/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jinzhu/gorm"
)

func helloWorld(c *fiber.Ctx) error {
	// return c.SendString("Hello, World!")
	return c.SendString(fmt.Sprintf("%d", time.Now().Unix()))
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "info.db")
	if err != nil {
		panic("Failed to connect to the Database")
	}

	fmt.Println("Database connection successfully established")

	database.DBConn.AutoMigrate(&controllers.Process{})
	database.DBConn.AutoMigrate(&controllers.Comment{})
	// database.DBConn.AutoMigrate(&controllers.User{})
	database.DBConn.AutoMigrate(&controllers.UserFile{})
	fmt.Println("DB auto-migration was set up")
}

func setupRouter(app *fiber.App) {
	app.Get("/api/v1/", helloWorld)
	app.Post("/api/v1/process", controllers.NewProcess)

	app.Post("/api/v1/comment", controllers.NewComment)
	app.Put("/api/v1/comment/:id", controllers.UpdateComment)
	app.Delete("/api/v1/comment/:id", controllers.DeleteComment)
	app.Get("/api/v1/comment/:id", controllers.GetComment)
	app.Get("/api/v1/comments", controllers.GetComments)
	app.Get("/api/v1/comments/process/:id", controllers.GetCommentsByProcessID)

	app.Post("/api/v1/files", controllers.NewFormFiles)
	app.Get("/api/v1/files", controllers.GetFilesWithoutBlob)
	app.Get("/api/v1/file/:id", controllers.GetFile)
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		BodyLimit:     1024 * 1024 * 1024,
		// StrictRouting: true,
		// ServerHeader:  "Fiber",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	initDatabase()

	// This is only executed at the end of the program,
	// for closing the DB connections
	defer database.DBConn.Close()

	setupRouter(app)

	app.Static("/", "./public")

	app.Listen(":3123")
}
