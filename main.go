package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/lucashmorais/documenta/controllers"
	"github.com/lucashmorais/documenta/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
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
	new_auth_function := jwtware.New(jwtware.Config{
		ErrorHandler: errorHandler,
		SigningKey:   []byte(jwtSecret),
		Filter:       authorizationFilter,
	})

	final_function := func(ctx *fiber.Ctx) error {
		ctx.Set("Cache-Control", "no-store, must-revalidate")
		ctx.Set("Pragma", "no-cache")
		ctx.Set("Expires", "0")
		return new_auth_function(ctx)
	}

	return final_function
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

func testTokenParsing(ctx *fiber.Ctx) error {
	// Token from another example.  This token is expired
	tokenString := ctx.Cookies("documentaLoginToken")

	if tokenString == "" {
		return ctx.Next()
	}

	fmt.Println("[testTokenParsing::tokenString]: " + tokenString)

	// The function that is passed as the second argument to jwt.Parse should, once called, return a byte slice containing the key to validate the token.
	// Such second argument is a function rather than a string or a byte array because the key might have to be dynamically determined based on information
	// about the token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if token.Valid {
		fmt.Println("[testTokenParsing]: You look nice today")

		fmt.Printf("[testTokenParsing::Locals::userID]: %v\n", controllers.RetrieveUserID(ctx))

	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("[testTokenParsing]: That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("[testTokenParsing]: Timing is everything")
		} else {
			fmt.Println("[testTokenParsing]: Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("[testTokenParsing]: Couldn't handle this token:", err)
	}

	// Without this, the next middleware or handler will never be called
	return ctx.Next()
}

func helloWorld(c *fiber.Ctx) error {
	// return c.SendString("Hello, World!")
	return c.SendString(fmt.Sprintf("%d", time.Now().Unix()))
}

func initDatabase() {
	var err error
	server_path := getServerPath()
	database.DBConn, err = gorm.Open("sqlite3", server_path+"/info.db")
	if err != nil {
		panic("Failed to connect to the Database")
	}

	fmt.Println("Database connection successfully established")

	database.DBConn.AutoMigrate(&controllers.Process{})
	database.DBConn.AutoMigrate(&controllers.Center{})
	database.DBConn.AutoMigrate(&controllers.ProcessType{})
	database.DBConn.AutoMigrate(&controllers.ProcessStatus{})
	database.DBConn.AutoMigrate(&controllers.Minute{})
	database.DBConn.AutoMigrate(&controllers.MinuteVersion{})

	database.DBConn.AutoMigrate(&controllers.Comment{})
	database.DBConn.AutoMigrate(&controllers.User{})
	database.DBConn.AutoMigrate(&controllers.UserSequence{})
	database.DBConn.AutoMigrate(&controllers.UserSequenceKind{})
	database.DBConn.AutoMigrate(&controllers.UserSequenceRank{})
	database.DBConn.AutoMigrate(&controllers.TokenPassingTimestamp{})
	database.DBConn.AutoMigrate(&controllers.Role{})
	database.DBConn.AutoMigrate(&controllers.Permission{})
	database.DBConn.AutoMigrate(&controllers.UserFile{})
	fmt.Println("DB auto-migration was set up")
}

func setupRouter(app *fiber.App) {
	app.Get("/api/v1/", helloWorld)
	// app.Get("/api/v1/login", controllers.Login)
	app.Post("/api/v1/login", controllers.Login)
	app.Get("/api/v1/users", controllers.GetUsers)

	protected := app.Group("/api/v1", authRequired())

	protected.Get("processes", controllers.GetProcesses)
	protected.Post("process", controllers.PostProcess)
	protected.Put("process/:process_id", controllers.PutProcess)
	protected.Put("process/:process_id/set_status/:status_id", controllers.PatchProcessStatus)

	protected.Get("process_types", controllers.GetProcessTypes)

	protected.Get("centers", controllers.GetCenters)

	protected.Post("comment", controllers.NewComment)
	// protected.Put("comment/:id", controllers.UpdateComment)
	protected.Put("comment", controllers.UpdateComment)
	protected.Delete("comment/:id", controllers.DeleteComment)
	protected.Get("comment/:id", controllers.GetComment)
	protected.Get("comments", controllers.GetComments)

	protected.Post("files", controllers.NewFormFiles)
	protected.Get("files", controllers.GetFilesWithoutBlob)
	protected.Get("file/:id", controllers.GetFile)
	protected.Delete("file/:id", controllers.DeleteFile)

	protected.Get("minutes", controllers.GetMinutes)
	protected.Post("minute", controllers.NewMinute)
	protected.Delete("minute/:id", controllers.DeleteMinute)
	protected.Patch("minute/:id", controllers.PatchMinute)

	protected.Get("minute_versions", controllers.GetMinuteVersions)
	protected.Post("minute_version", controllers.NewMinuteVersion)
	protected.Delete("minute_version/:id", controllers.DeleteMinuteVersion)

	protected.Get("current_user", controllers.GetLoggedUser)
	protected.Get("current_user/permissions", controllers.GetLoggedUserPermissions)
	protected.Post("user", controllers.PostUser)
	protected.Put("user", controllers.PutUser)
	protected.Delete("users", controllers.DeleteUsers)

	protected.Get("user_sequences", controllers.GetUserSequences)
	protected.Get("user_sequence_with_timestamps", controllers.GetLastUserSequenceForGivenProcessWithTimestamps)
	protected.Get("token_passing_timestamps", controllers.GetTokenPassingTimestamps)
	protected.Post("user_sequence", controllers.PostUserSequence)
	protected.Post("user_sequence_simple", controllers.PostUserSequenceSimple)
	protected.Post("user_sequence/count_completion", controllers.IncreaseSequenceCompletionCounter)

	protected.Get("roles", controllers.GetRoles)
	protected.Post("role", controllers.PostRole)
	protected.Put("role", controllers.PutRole)
	protected.Delete("roles", controllers.DeleteRoles)

	protected.Get("permissions", controllers.GetPermissions)

	protected.Get("next_inbound_protocol_number", controllers.GetNextInboundProtocolNumber)
	protected.Get("next_outbound_protocol_number", controllers.GetNextOutboundProtocolNumber)
	protected.Patch("minute_protocol_numbers/:id", controllers.PatchMinuteProtocol)
}

func addAuthRequestHeader(ctx *fiber.Ctx) error {
	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIxLTA2LTIzVDE4OjE2OjQ5LjcxMDU1Mzg3Ny0wMzowMCIsInN1YiI6IjEifQ.T7CADK7tFePIi_d8lcw4PS5RMLFIBu51j_rmdoHaDd8"
	token := ctx.Cookies("documentaLoginToken")
	// println(token)
	// ctx.Context().Request.Header.Add("Authorization", "Bearer "+token)
	ctx.Context().Request.Header.Add("Authorization", "Bearer "+token)
	ctx.Context().Request.Header.Add("Cache-Control", "no-cache, must-revalidate")
	ctx.Context().Request.Header.Add("Pragma", "no-cache")
	// println(string(ctx.Context().Request.Header.Header()))
	ctx.Next()
	return nil
}

func getServerPath() string {
	env_path, env_path_is_set := os.LookupEnv("DOCUMENTA_ROOT")
	if env_path_is_set {
		return env_path
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		return exPath
	}
}

func main() {
	server_path := getServerPath()
	fmt.Printf("Server path: %s\n", server_path)
	godotenv.Load(server_path + "/config/config.env")

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		BodyLimit:     1024 * 1024 * 1024,
		// StrictRouting: true,
		// ServerHeader:  "Fiber",
	})

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(addAuthRequestHeader)
	// app.Use("document.html", authRequired())
	// app.Use("index.html", authRequired())
	app.Use("/document.html", authRequired())
	// app.Use("/home.html", authRequired(), testTokenParsing)
	app.Use("/home.html", authRequired())
	app.Use("/register.html", authRequired())
	// app.Use("/", logJWTInformation)
	// app.Use("/", testTokenParsing)

	initDatabase()

	// This is only executed at the end of the program,
	// for closing the DB connections
	defer database.DBConn.Close()

	setupRouter(app)

	app.Static("/", server_path+"/public")

	app.Listen(":3123")
}
