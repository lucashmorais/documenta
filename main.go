package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lucashmorais/documenta/controllers"
	"github.com/lucashmorais/documenta/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/jinzhu/gorm"
)

const jwtSecret = "asecret"

func errorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.Redirect("/")
}

func authorizationFilter(ctx *fiber.Ctx) bool {
	return ctx.Path() == "/"
}

func authRequired() func(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		ErrorHandler: errorHandler,
		SigningKey:   []byte(jwtSecret),
		Filter:       authorizationFilter,
	})
}

func logJWTInformation(ctx *fiber.Ctx) error {
	user_token := ctx.Locals("user").(*jwt.Token)
	// user_token := ctx.Locals("user")
	// fmt.Printf("%T\n", user_token)
	// println(user_token)
	token_claims := user_token.Claims.(jwt.MapClaims)
	id_from_claims := token_claims["sub"].(string)
	expiration := token_claims["exp"].(string)
	fmt.Printf("Interaction from JWT-authorized user with id %s is allowed up to %s\n", id_from_claims, expiration)
	return ctx.Next()
}

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
	database.DBConn.AutoMigrate(&controllers.User{})
	database.DBConn.AutoMigrate(&controllers.UserFile{})
	fmt.Println("DB auto-migration was set up")
}

func setupRouter(app *fiber.App) {
	app.Get("/api/v1/", helloWorld)
	// app.Get("/api/v1/login", controllers.Login)
	app.Post("/api/v1/login", controllers.Login)
	app.Get("/api/v1/users", controllers.GetUsers)

	protected := app.Group("/api/v1", authRequired())
	protected.Post("process", controllers.NewProcess)

	protected.Post("comment", controllers.NewComment)
	protected.Put("comment/:id", controllers.UpdateComment)
	protected.Delete("comment/:id", controllers.DeleteComment)
	protected.Get("comment/:id", controllers.GetComment)
	protected.Get("comments", controllers.GetComments)
	protected.Get("comments/process/:id", controllers.GetCommentsByProcessID)

	protected.Post("files", controllers.NewFormFiles)
	protected.Get("files", controllers.GetFilesWithoutBlob)
	protected.Get("file/:id", controllers.GetFile)
	protected.Delete("file/:id", controllers.DeleteFile)

	protected.Post("user", controllers.PostUser)
}

func addAuthRequestHeader(ctx *fiber.Ctx) error {
	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIxLTA2LTIzVDE4OjE2OjQ5LjcxMDU1Mzg3Ny0wMzowMCIsInN1YiI6IjEifQ.T7CADK7tFePIi_d8lcw4PS5RMLFIBu51j_rmdoHaDd8"
	token := ctx.Cookies("documentaLoginToken")
	println(token)
	// ctx.Context().Request.Header.Add("Authorization", "Bearer "+token)
	ctx.Context().Request.Header.Add("Authorization", "Bearer "+token)
	// println(string(ctx.Context().Request.Header.Header()))
	ctx.Next()
	return nil
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
	app.Use(addAuthRequestHeader)
	// app.Use("document.html", authRequired())
	// app.Use("index.html", authRequired())
	app.Use("/document.html", authRequired())
	app.Use("/home.html", authRequired())
	app.Use("/register.html", authRequired())
	// app.Use("/", logJWTInformation)

	initDatabase()

	// This is only executed at the end of the program,
	// for closing the DB connections
	defer database.DBConn.Close()

	setupRouter(app)

	app.Static("/", "./public")

	app.Listen(":3123")
}
