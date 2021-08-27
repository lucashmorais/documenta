package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

// Struct that defines fields of the Minute gorm table
type Minute struct {
	gorm.Model
	Content     string
	Description string
	UserID      int
	User        User
	ProcessID   int
	Process     Process
	Center      Center
	CenterID    int

	//TODO: ENSURE THAT BOTH OF THESE ARE UNIQUE!
	InboundProtocol  string
	OutboundProtocol string

	UnixCreatedAt int64
	UnixUpdatedAt int64
	UnixDeletedAt int64
}

// Struct that defines fields of the MinuteVersions gorm table
type MinuteVersion struct {
	gorm.Model
	Content     string
	Description string
	UserID      int
	User        User
	Minute      Minute
	MinuteID    int

	UnixCreatedAt int64
	UnixUpdatedAt int64
	UnixDeletedAt int64
}

// Function that fetches all items in the Minute table from te gorm database and returns them as a JSON array
// If `processID` is provided, it only returns Minutes related to that Process
// If `centerID` is provided, it only returns Minutes related to that Center
// If `userID` is provided, it only returns Minutes related to that User
func GetMinutes(c *fiber.Ctx) error {
	db := database.DBConn
	var minutes []Minute

	driver := db.Preload("User").Preload("Process").Preload("Center")

	processID := c.Query("processID")
	centerID := c.Query("centerID")
	userID := c.Query("userID")

	if processID != "" {
		i, err := strconv.Atoi(processID)
		if err == nil {
			driver = driver.Where("process_id = ?", i)
		}
	}

	if centerID != "" {
		i, err := strconv.Atoi(centerID)
		if err == nil {
			driver = driver.Where("center_id = ?", i)
		}
	}

	if userID != "" {
		i, err := strconv.Atoi(userID)
		if err == nil {
			driver = driver.Where("user_id = ?", i)
		}
	}

	driver.Find(&minutes)

	return c.JSON(minutes)
}

// Function that adds a new item to the Minute gorm table based on the data provided in the POST request as a JSON
func NewMinute(c *fiber.Ctx) error {
	db := database.DBConn
	// var process Minute

	minute := Minute{}

	err := c.BodyParser(&minute)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	userID := RetrieveUserID(c)

	if userID == 0 {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "could_not_retrieve_user_id",
		})
	}

	// fmt.Printf("[PostMinute]: Decoded new minute: %v\n", minute)

	minute.UnixCreatedAt = time.Now().Unix()
	minute.UserID = userID
	db.Create(&minute)

	// fmt.Printf("[PostMinute]: Minute with userID: %v\n", minute)

	return c.JSON(minute)
}

// Function that adds a new item to the MinuteVersion gorm table based on the data provided in the POST request as a JSON
func NewMinuteVersion(c *fiber.Ctx) error {
	db := database.DBConn

	minuteVersion := MinuteVersion{}

	err := c.BodyParser(&minuteVersion)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	fmt.Printf("[PostMinuteVersion]: Decoded new minute version: %v\n", minuteVersion)

	db.Create(&minuteVersion)

	return c.JSON(minuteVersion)
}

// Function that fetches all items in the MinuteVersion table from te gorm database and returns them as a JSON array
// If `minuteID` is provided, it only returns MinuteVersions related to that Minute
func GetMinuteVersions(c *fiber.Ctx) error {
	db := database.DBConn
	var minutes []MinuteVersion

	minuteID := c.Query("minuteID")

	driver := db.Preload("User")

	if minuteID != "" {
		i, err := strconv.Atoi(minuteID)
		if err == nil {
			driver = driver.Where("minute_id = ?", i)
		}
	}

	driver.Find(&minutes)

	return c.JSON(minutes)
}

// Function that deletes a Minute based on the ID provided in the DELETE request URL path
func DeleteMinute(c *fiber.Ctx) error {
	db := database.DBConn

	minuteID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "invalid_id",
		})
	}

	minute := Minute{}

	db.First(&minute, minuteID)

	if minute.ID == 0 {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"cause":   "not_found",
		})
	}

	db.Delete(&minute)

	return c.JSON(&fiber.Map{
		"success": true,
	})
}

// Function that deletes a MinuteVersion based on the ID provided in the DELETE request URL path
func DeleteMinuteVersion(c *fiber.Ctx) error {
	db := database.DBConn

	minuteVersionID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "invalid_id",
		})
	}

	minuteVersion := MinuteVersion{}

	db.First(&minuteVersion, minuteVersionID)

	if minuteVersion.ID == 0 {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"cause":   "not_found",
		})
	}

	db.Delete(&minuteVersion)

	return c.JSON(&fiber.Map{
		"success": true,
	})
}
