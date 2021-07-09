package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

const jwtSecret = "asecret"

func Login(ctx *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ctx.JSON(fiber.Map{"derp": "derp"})

	var body request
	body.Email = "bob@gmail.com"
	body.Password = "password123"

	err := ctx.BodyParser(&body)

	// var err error = nil
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return fiber.ErrBadRequest
	}

	if body.Email != "bob@gmail.com" || body.Password != "password123" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
		return fiber.ErrBadRequest
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Subject stands for "subject", and is basically a user identifier.
	// If emails are unique across the whole system, they could be used
	// as valid subject claims.
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7) // a week

	s, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": s,
		"user": struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		}{
			Id:    1,
			Email: "bob@gmail.com",
		},
	})

	return nil
}
