package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hlandau/passlib"
	"github.com/lucashmorais/documenta/database"
)

const jwtSecret = "asecret"

func Login(ctx *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request

	err := ctx.BodyParser(&body)

	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return fiber.ErrBadRequest
	}

	db := database.DBConn
	var user User

	db.Where("email = ?", body.Email).Find(&user)
	_, err = passlib.Verify(body.Password, user.PHash)

	if err != nil {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
		return fiber.ErrBadRequest
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	expiration := time.Now().Add(time.Hour * 24 * 7) // a week

	// Subject stands for "subject", and is basically a user identifier.
	// If emails are unique across the whole system, they could be used
	// as valid subject claims.
	claims["sub"] = "1"
	claims["exp"] = expiration

	s, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "documentaLoginToken",
		Value:   s,
		Expires: expiration,
	})

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

func GetUser(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var user User
	db.Where("email = ?", "bob@gmail.com").Where("p_hash = ?", "password123").Find(&user)
	// db.Where("Name = ?", "Albert Billford").Find(&user)

	return c.JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var user []User
	db.Find(&user)
	// db.Where("Name = ?", "Albert Billford").Find(&user)

	return c.JSON(user)
}

func PostUser(c *fiber.Ctx) error {
	db := database.DBConn
	var user, oldUser User
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	db.Where("email = ?", user.Email).Find(&oldUser)

	if oldUser.Email != "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "email_was_taken",
		})
	}

	hash, err := passlib.Hash(user.PHash)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "cannot_hash_password",
		})
	}

	fmt.Printf("[hash from user-provided password]: %s", hash)

	user.PHash = hash

	db.Create(&user)
	return c.JSON(user)
}
