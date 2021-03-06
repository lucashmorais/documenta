package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hlandau/passlib"
	"github.com/lucashmorais/documenta/database"
)

// TODO: CHANGE BEFORE DEPLOYMENT!
const jwtSecret = "asecret"

func RetrieveUserID(ctx *fiber.Ctx) int {
	tokenString := ctx.Cookies("documentaLoginToken")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	token_claims := token.Claims.(jwt.MapClaims)
	i := token_claims["user-id"]

	if err == nil {
		return int(i.(float64))
	}

	return 0
}

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
	claims["user-email"] = user.Email
	claims["user-id"] = user.ID
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
