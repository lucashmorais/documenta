package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type Process struct {
	gorm.Model
	Title     string `json: "title"`
	Summary   string `json: "summary"`
	Reference int    `json: "reference"`

	// These are all "belongs to" relationships
	CenterID        uint
	Center          Center
	UserID          uint
	User            User
	ProcessStatusID uint
	ProcessStatus   ProcessStatus
	ProcessTypeID   uint
	ProcessType     ProcessType
	UserSequence    UserSequence
}

type ProcessStatus struct {
	gorm.Model
	Name        string
	Description string
}

func GetProcesses(c *fiber.Ctx) error {
	db := database.DBConn
	var processes []Process

	processIDRaw := c.Query("processID")
	statusString := c.Query("statusString")
	statusIDRaw := c.Query("statusID")
	typeString := c.Query("typeString")
	typeIDRaw := c.Query("typeID")
	onlyModifiableByUser := c.Query("onlyModifiableByUser")

	// IMPORTANT: The following query will only work for processes that have valid references
	// to ProcessStatus and ProcessTypes, since it includes joins with these tables
	base := db.
		Preload("Center").
		Preload("ProcessStatus").
		Preload("ProcessType").
		Table("processes").
		Joins("join process_statuses on process_statuses.id = processes.process_status_id").
		Joins("join process_types on process_types.id = processes.process_type_id")

	if statusString != "" {
		base = base.Where("process_statuses.name = ?", statusString)
	}
	if typeString != "" {
		base = base.Where("process_types.name = ?", typeString)
	}
	if processIDRaw != "" {
		i, err := strconv.Atoi(processIDRaw)
		if err == nil {
			base = base.Where("processes.id = ?", i)
		}
	}
	if statusIDRaw != "" {
		i, err := strconv.Atoi(statusIDRaw)
		if err == nil {
			base = base.Where("process_statuses.id = ?", i)
		}
	}
	if typeIDRaw != "" {
		i, err := strconv.Atoi(typeIDRaw)
		if err == nil {
			base = base.Where("process_types.id = ?", i)
		}
	}

	base.Find(&processes)

	// If `onlyModifiableUser == "true"`, the following piece of code
	// builds a new Process slice from `processes` such that, for each and all of its elements,
	// the `num_completions` field of the most recent `UserSequence` related to it is equal to
	// the position of the the current logged user in the `user_sequence_users` slice whose all
	// elements are related to that `UserSequence`.
	//
	// (interesting) NOTE: Github copilot failed to implement this
	if onlyModifiableByUser == "true" {
		userID := RetrieveUserID(c)
		processesTheCurrentUserCanModify := []Process{}

		for idx, proc := range processes {
			var userSeq UserSequence
			db.Where("process_id = ?", proc.ID).Last(&userSeq)

			if userSeq.ID != 0 {
				var users []User

				// This finds users belonging to the many2many relationship
				// between UserSequence and Users
				db.Model(&userSeq).Association("Users").Find(&users)

				if len(users) > userSeq.NumCompletions {
					activeUser := users[userSeq.NumCompletions]
					if userID == int(activeUser.ID) {
						processesTheCurrentUserCanModify = append(processesTheCurrentUserCanModify, processes[idx])
					}
				}
			}
		}

		processes = processesTheCurrentUserCanModify
	}

	return c.JSON(processes)
}

func PostProcess(c *fiber.Ctx) error {
	db := database.DBConn
	// var process Process

	process := struct {
		gorm.Model
		Title               string `json: "title"`
		Summary             string `json: "summary"`
		TypeID              uint   `json: "typeID"`
		CenterID            uint   `json: "centerID"`
		UserSequenceUserIDs []uint `json: "userSequenceUserIDs"`
	}{}

	err := c.BodyParser(&process)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	fmt.Printf("[PostProcess]: Decoded new process: %v\n", process)

	var processType ProcessType
	db.Where(process.TypeID).Find(&processType)

	var center Center
	db.Where(process.CenterID).Find(&center)

	var status ProcessStatus
	db.Where("name = ?", "Rascunho").Find(&status)

	// This retrieves a slice of User objects matching a slice of UserIDs
	var users []User

	// WARNING: THE FOLLOWING DOES NOT WORK, BECAUSE IT ORDERS USERS BY ID,
	// DISREGARDING THE ORDER OF THE USERIDS IN THE USERSEQUENCEUSERIDS SLICE
	//
	// db.Where("id IN (?)", process.UserSequenceUserIDs).Find(&users)

	for _, userID := range process.UserSequenceUserIDs {
		var user User
		db.Where(userID).Find(&user)
		users = append(users, user)
	}

	userSequence := UserSequence{
		Users:              users,
		NumUsers:           len(users),
		NumCompletions:     0,
		UserSequenceKindID: int(Revision),
	}

	db.Create(&userSequence)

	dbProcess := Process{Title: process.Title, Summary: process.Summary, ProcessTypeID: process.TypeID, ProcessType: processType, Center: center, ProcessStatus: status, UserSequence: userSequence}

	db.Create(&dbProcess)

	// Add information about the order at which Users have been included in the UserSequenceUserIDs
	// to the UserSequenceRanks table
	for idx, user := range users {
		userSequenceRank := UserSequenceRank{
			UserSequenceID: int(userSequence.ID),
			UserID:         int(user.ID),
			Rank:           idx,
		}

		db.Create(&userSequenceRank)
	}

	return c.JSON(dbProcess)
}

func PutProcess(c *fiber.Ctx) error {
	db := database.DBConn
	// var process Process

	process := struct {
		gorm.Model
		Title               string `json: "title"`
		Summary             string `json: "summary"`
		TypeID              uint   `json: "typeID"`
		CenterID            uint   `json: "centerID"`
		UserSequenceUserIDs []uint `json: "userSequenceUserIDs"`
	}{}

	err := c.BodyParser(&process)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	fmt.Printf("[PutProcess]: Decoded new process: %v\n", process)

	var processType ProcessType
	db.Where(process.TypeID).Find(&processType)

	var center Center
	db.Where(process.CenterID).Find(&center)

	// This retrieves a slice of User objects matching a slice of UserIDs
	var users []User

	// WARNING: THE FOLLOWING DOES NOT WORK, BECAUSE IT ORDERS USERS BY ID,
	// DISREGARDING THE ORDER OF THE USERIDS IN THE USERSEQUENCEUSERIDS SLICE
	//
	// db.Where("id IN (?)", process.UserSequenceUserIDs).Find(&users)

	for _, userID := range process.UserSequenceUserIDs {
		var user User
		db.Where(userID).Find(&user)
		users = append(users, user)
	}

	var userSequence UserSequence

	// Get current user sequence from DB if the `keep_user_sequence` query param is set to true.
	// Otherwise, create a new user sequence using the `users` slice.
	if c.Query("keep_user_sequence") == "true" {
		db.Where("process_id = ?", c.Params("id")).Find(&userSequence)
	} else {
		userSequence = UserSequence{
			Users:              users,
			NumUsers:           len(users),
			NumCompletions:     0,
			UserSequenceKindID: int(Revision),
		}
	}

	var dbProcess Process
	db.Where("id = ?", c.Params("process_id")).Find(&dbProcess)

	dbProcess.Title = process.Title
	dbProcess.Summary = process.Summary
	dbProcess.ProcessTypeID = process.TypeID
	dbProcess.ProcessType = processType
	dbProcess.Center = center
	dbProcess.UserSequence = userSequence

	fmt.Printf("[PutProcess::dbProcess]: %v\n", dbProcess)

	db.Save(&dbProcess)

	if c.Query("keep_user_sequence") != "true" {
		// Add information about the order at which Users have been included in the UserSequenceUserIDs
		// to the UserSequenceRanks table
		for idx, user := range users {
			userSequenceRank := UserSequenceRank{
				UserSequenceID: int(dbProcess.UserSequence.ID),
				UserID:         int(user.ID),
				Rank:           idx,
			}

			db.Create(&userSequenceRank)
		}
	}

	return c.JSON(dbProcess)
}

// Function that updates a `Process`'s `ProcessStatus`
func PatchProcessStatus(c *fiber.Ctx) error {
	db := database.DBConn

	processID := c.Params("process_id")
	statusID := c.Params("status_id")

	var process Process
	db.Where("id = ?", processID).Find(&process)

	var status ProcessStatus
	db.Where("id = ?", statusID).Find(&status)

	process.ProcessStatus = status

	db.Save(&process)

	return c.JSON(process)
}
