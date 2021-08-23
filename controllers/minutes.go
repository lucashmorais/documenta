package controllers

import (
	"strconv"

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
}

// Struct that defines fields of the MinuteVersions gorm table
type MinuteVersion struct {
	gorm.Model
	Content     string
	Description string

	UnixCreatedAt int64
	UnixUpdatedAt int64
	UnixDeletedAt int64
}

// Function that fetches all items in the Minute table from te gorm database and returns them as a JSON array
func GetMinutes(c *fiber.Ctx) error {
	db := database.DBConn
	var minutes []Minute

	driver := db.Preload("User").Preload("Process")
	driver.Find(&minutes)

	return c.JSON(minutes)
}

// Function that fetches all items in the MinuteVersion table from te gorm database and returns them as a JSON array
// If `minuteID` is provided, it only returns MinuteVersions related to that Minute
func GetMinuteVersions(c *fiber.Ctx) error {
	db := database.DBConn
	var minutes []MinuteVersion

	minuteID := c.Query("processID")

	driver := db.Preload("User").Preload("Process")

	if minuteID != "" {
		i, err := strconv.Atoi(minuteID)
		if err == nil {
			driver = driver.Where("minute_id = ?", i)
		}
	}

	driver.Find(&minutes)

	return c.JSON(minutes)
}
