package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	jwtware "github.com/gofiber/jwt"
)

const jwtSecret = "asecret"

// This function build another one that merely checks whether the user has provided a valid
// JWT issued by this server. If the user has one, it is because he
// has already logged in to the server.
//
// A user will only be able to provide a valid JWT insofar as he or she
// was able to successfully login with certain user credentials. The
// JWT thus obtained will thus let the user interact with this server
// with the permissions that have been ascribed to the user record related
// to that JWT.
//
// This function follows the factory pattern.
func authRequired() func(ctx *fiber.Ctx) {
	return jwtware.New(jwtware.Config{
		// The function to which the `ErrorHandler` struct field
		// is initialized is called whenever the JWT library detects
		// the JWT coming to the server to be invalid, for any reason.
		ErrorHandler: func(ctx *fiber.Ctx, err error) {
			ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
		// The `SigningKey` struct field indicates which encryption key
		// should be used for decoding user information.
		//
		// The existence and secrecy of this key ensures that
		// if this server is able to decode a JWT with such key, the
		// JWT had almost certainly been generated and sent to the user
		// by this server or some party with which it had explicitly
		// shared its signing key.
		SigningKey: []byte(jwtSecret),
	})
}

func main() {
	app := fiber.New()

	app.Use(middleware.Logger())

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("Hello world")
	})

	app.Post("/login", login)

	// The user will only succeed with this GET request if he or she
	// is able to provide a valid JWT to `authRequired()`, the middleware
	// that is executed prior to the fulfilling of the request.
	//
	// By the time the server starts running the anonymous function passed
	// as the handler for this route, the user will have already provided
	// a JWT related to some valid user recognized by the backend.
	app.Get("/hello", authRequired(), func(ctx *fiber.Ctx) {
		// Q1: Where is this "user" local data coming from?
		// H1: It has probably been initialized by the jwtware constructor
		// Q3: What information is contained in this JWT?
		// A3: All that was encoded by the server as it built
		//     the JWT to be passed to the user as a result
		//     of the latter's success at logging in to the server.
		//
		//     In particular, it includes a simple ID that might
		//     be used for pooling a DB to retrieve information
		//     about such user.
		// Q4: What is the `.(*jwt.Token)` construct used for?
		//     for casting the interface{} value returned by the Locals
		//     method to the `* jwt.Token` type.
		//
		//     In fact, the two following lines also use the same
		//     kind of construct to cast `interface{}` values to more
		//     specific types of variables.
		user_token := ctx.Locals("user").(*jwt.Token)
		token_claims := user_token.Claims.(jwt.MapClaims)
		id_from_claims := token_claims["sub"].(string)
		ctx.Send(fmt.Sprintf("Hello user with id: %s", id_from_claims))
	})

	err := app.Listen(3000)
	if err != nil {
		panic(err)
	}
}

func login(ctx *fiber.Ctx) {
	// Hard-coded login information
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	err := ctx.BodyParser(&body)

	// Reply with error status if login information
	// cannot be decoded as a JSON
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return
	}

	// Reply with a bad credentials status if login information
	// does not match data for server authorized users
	if body.Email != "bob@gmail.com" || body.Password != "password123" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
		return
	}

	// IMPORTANT: Initialize new JWT
	token := jwt.New(jwt.SigningMethodHS256)

	// IMPORTANT: Initialize new claims object
	//
	//            Claims objects are used for storing any information
	//            that might be apposite for identifying and serving
	//            a user holding such JWT in the future.
	//            At a manimum level, such objects should provide some
	//            unique identifier to the user that might be used as
	//            a primary key to access further information about such
	//            user at some central user database.
	claims := token.Claims.(jwt.MapClaims)
	// Q2: What is this "sub" claim?
	// A2: Merely a numerical ID for a certain user.
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7) // a week

	// This function call transforms the Golang `token` object into `s`, a
	// a signed string that might by passed without further transformations
	// to the user.
	s, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
		return
	}

	// This sends a JSON back to the user containing the signed
	// JWT string that the latter might use to interact with the
	// system as a logged agent. Secondarily, it also contains
	// other unencoded information that the server internally
	// uses to identify such user.
	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		// The `token` JSON field
		"token": s,
		"user": struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		}{
			Id:    1,
			Email: "bob@gmail.com",
		},
	})
}
