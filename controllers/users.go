package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hlandau/passlib"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type User struct {
	gorm.Model
	Name      string `json: "name"`
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
	Title     string `json: "title"`
	Initials  string `json: "initials"`
	Email     string `json: "email"`
	PHash     string `json: "phash"`
	Roles     []Role `gorm:"many2many:user_roles"`
}

type UserSequence struct {
	gorm.Model
	ProcessID          int
	NumUsers           int
	NumCompletions     int
	UserSequenceKindID int
	UserSequenceKind   UserSequenceKind
	Users              []User `gorm:"many2many:user_sequence_users"`
}

type UserSequenceKind struct {
	ID   int
	Name string
}

type UserSequenceRank struct {
	gorm.Model
	UserSequenceID int
	UserID         int
	Rank           int
}

type UserSequenceKindIDValue int

const (
	Revision UserSequenceKindIDValue = 1
	Approval UserSequenceKindIDValue = 2
)

type TokenPassingTimestamp struct {
	gorm.Model
	UserSequenceID int
	UserSequence   UserSequence
	UnixTimestamp  int64
}

func GetUser(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var user User
	db.Preload("Roles").Where("email = ?", "bob@gmail.com").Where("p_hash = ?", "password123").Find(&user)
	// db.Where("Name = ?", "Albert Billford").Find(&user)

	return c.JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	// email := c.Params("email")
	// password := c.Params("password")

	db := database.DBConn
	var user []User
	db.Preload("Roles").Find(&user)
	// db.Where("Name = ?", "Albert Billford").Find(&user)

	return c.JSON(user)
}

func PostUser(c *fiber.Ctx) error {
	db := database.DBConn
	var oldUser User

	user := struct {
		gorm.Model
		PHash     string `json: "phash"`
		Email     string `json: "email"`
		FirstName string `json: "firstName"`
		LastName  string `json: "lastName"`
		Initials  string `json: "initials"`
		Roles     []int  `json: "roles"`
	}{}

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

	// fmt.Printf("[hash from user-provided password]: %s\n", hash)

	user.PHash = hash

	var roleSlice []Role
	db.Find(&roleSlice, user.Roles)
	fmt.Println(roleSlice)

	dbUser := User{PHash: user.PHash, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName, Initials: user.Initials, Roles: roleSlice}
	db.Create(&dbUser)

	return c.JSON(dbUser)
}

func PutUser(c *fiber.Ctx) error {
	db := database.DBConn
	user := struct {
		gorm.Model
		Email     string `json: "email"`
		FirstName string `json: "firstName"`
		LastName  string `json: "lastName"`
		Initials  string `json: "initials"`
		Roles     []int  `json: "roles"`
	}{}

	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	dbUser := User{}

	var roleSlice []Role

	db.Where("id = ?", user.ID).First(&dbUser)
	db.Find(&roleSlice, user.Roles)

	db.Model(&dbUser).Update(user)
	db.Model(&dbUser).Association("Roles").Replace(roleSlice)
	return c.JSON(user)
}

func DeleteUsers(c *fiber.Ctx) error {
	db := database.DBConn
	idsToDelete := struct {
		gorm.Model
		Ids []int `json: "ids"`
	}{}

	err := c.BodyParser(&idsToDelete)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	fmt.Println("[idsToDelete]", idsToDelete)
	fmt.Println("[idsToDelete.Ids]", idsToDelete.Ids)
	db.Delete(&User{}, idsToDelete.Ids)

	return c.SendStatus(200)
}

func GetUserSequences(c *fiber.Ctx) error {
	db := database.DBConn
	var userSequences []UserSequence

	driver := db.Preload("Users").Preload("UserSequenceKind")
	processID := c.Query("processID")

	// TODO: Check whether we should convert processID to a number before this
	if processID != "" {
		driver = driver.Where("process_id = ?", processID)
	}

	driver.Find(&userSequences)

	return c.JSON(userSequences)
}

// Function that retrieves all TokenPassingTimestamps for a given UserSequence specified by the `userSequenceID` query parameter
func GetTokenPassingTimestamps(c *fiber.Ctx) error {
	db := database.DBConn
	var tokenPassingTimestamps []TokenPassingTimestamp

	userSequenceID := c.Query("userSequenceID")

	if userSequenceID == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "user_sequence_id_not_specified",
		})
	} else {
		db.Where("user_sequence_id = ?", userSequenceID).Find(&tokenPassingTimestamps)
	}

	return c.JSON(tokenPassingTimestamps)
}

// Function that retrieves as a JSON object both (1) the most recent UserSequence corresponding to
// a certain `processID` query parameter and (2) all TokenPassingTimestamps for that UserSequence
// NOTE: This function was entirely written by co-pilot
func GetLastUserSequenceForGivenProcessWithTimestamps(c *fiber.Ctx) error {
	db := database.DBConn
	var userSequence UserSequence
	var tokenPassingTimestamps []TokenPassingTimestamp

	processID := c.Query("processID")

	if processID == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "process_id_not_specified",
		})
	} else {
		db.Where("process_id = ?", processID).Last(&userSequence)
		db.Where("user_sequence_id = ?", userSequence.ID).Find(&tokenPassingTimestamps)

		// Get the Users associated with the UserSequence ordered by rank
		var userSequenceRanks []UserSequenceRank
		db.Where("user_sequence_id = ?", userSequence.ID).Order("rank").Find(&userSequenceRanks)
		var users []User
		for _, userSequenceRank := range userSequenceRanks {
			var user User
			db.Where("id = ?", userSequenceRank.UserID).First(&user)
			users = append(users, user)
		}

		userSequence.Users = users
	}

	return c.JSON(map[string]interface{}{
		"userSequence":           userSequence,
		"tokenPassingTimestamps": tokenPassingTimestamps,
	})
}

// Function that retrieves a JSON object describing the current logged User identified by the JWT passed to the server
func GetLoggedUser(c *fiber.Ctx) error {
	db := database.DBConn
	userID := RetrieveUserID(c)

	var loggedUser User

	db.Where("id = ?", userID).First(&loggedUser)

	return c.JSON(loggedUser)
}

func PostUserSequence(c *fiber.Ctx) error {
	db := database.DBConn
	var userSequence UserSequence

	err := c.BodyParser(&userSequence)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	db.Create(&userSequence)

	return c.JSON(userSequence)
}

func PostUserSequenceSimple(c *fiber.Ctx) error {
	db := database.DBConn
	basicUserSequence := struct {
		gorm.Model
		UserSequenceKindID  int    `json: "userSequenceKindID"`
		UserSequenceUserIDs []uint `json: "userSequenceUserIDs"`
		ProcessID           int    `json: "processID"`
	}{}

	err := c.BodyParser(&basicUserSequence)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	var users []User
	db.Where("id IN (?)", basicUserSequence.UserSequenceUserIDs).Find(&users)

	userSequence := UserSequence{
		Users:              users,
		NumUsers:           len(users),
		NumCompletions:     0,
		UserSequenceKindID: basicUserSequence.UserSequenceKindID,
		ProcessID:          basicUserSequence.ProcessID,
	}

	db.Create(&userSequence)

	for idx, user := range users {
		userSequenceRank := UserSequenceRank{
			UserSequenceID: int(userSequence.ID),
			UserID:         int(user.ID),
			Rank:           idx,
		}

		db.Create(&userSequenceRank)
	}

	return c.JSON(userSequence)
}

// Function that increases the `NumCompletions` field of the most recent
// `UserSequence` corresponding to a certain `processID` query parameter
func IncreaseSequenceCompletionCounter(c *fiber.Ctx) error {
	db := database.DBConn
	var userSequence UserSequence

	processID := c.Query("processID")

	if processID == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "process_id_not_specified",
		})
	} else {
		db.Preload("Users").Where("process_id = ?", processID).Last(&userSequence)
		db.Model(&userSequence).Update("NumCompletions", userSequence.NumCompletions+1)
	}

	var tokenPassingTimestamp TokenPassingTimestamp
	tokenPassingTimestamp.UserSequenceID = int(userSequence.ID)
	tokenPassingTimestamp.UnixTimestamp = time.Now().Unix()
	db.Create(&tokenPassingTimestamp)

	return c.JSON(userSequence)
}

// If `permissionID` is equal to zero, this function returns a slice of all
// Permissions held by the logged user. Otherwise, it returns a non-empty
// Permission slice if, and only if, the user holds a Permission with DB index
// equal to `permissionID`.
func CoreGetLoggedUserPermissions(c *fiber.Ctx, permissionID int) []Permission {
	db := database.DBConn
	loggedUserID := RetrieveUserID(c)
	loggedUser := User{}

	db.Where(loggedUserID).Find(&loggedUser)

	var roles []Role
	db.Model(&loggedUser).Association("Roles").Find(&roles)

	var permissions []Permission

	for _, role := range roles {
		var rolePermissions []Permission
		if permissionID != 0 {
			db.Where("permission_id = ?", permissionID).Model(role).Association("Permissions").Find(&permissions)
			if len(permissions) > 0 {
				return permissions
			}
		} else {
			db.Model(role).Association("Permissions").Find(&permissions)
		}
		permissions = append(permissions, rolePermissions...)
	}

	return permissions
}

// Function that retrieves a JSON object describing all the Permissions of the currently logged User.
func GetLoggedUserPermissions(c *fiber.Ctx) error {
	return c.JSON(CoreGetLoggedUserPermissions(c, 0))
}
