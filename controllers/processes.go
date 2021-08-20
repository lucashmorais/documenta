package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

func GetProcesses(c *fiber.Ctx) error {
	db := database.DBConn
	var processes []Process

	statusString := c.Query("statusString")
	statusIDRaw := c.Query("statusID")

	if statusString != "" {
		var processStatus ProcessStatus
		db.Where("name = ?", statusString).Find(&processStatus)
		db.Preload("ProcessStatus").Preload("ProcessType").Where("process_status_id = ?", processStatus.ID).Find(&processes)
	} else if statusIDRaw != "" {
		statusID, err := strconv.Atoi(statusIDRaw)
		if err == nil {
			db.Preload("Center").Preload("ProcessStatus").Preload("ProcessType").Where("process_status_id = ?", statusID).Find(&processes)
		} else {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"cause":   "bad statusID query",
			})
		}
	} else {
		db.Preload("Center").Preload("ProcessStatus").Preload("ProcessType").Find(&processes)
	}

	return c.JSON(processes)
}

func PostProcess(c *fiber.Ctx) error {
	db := database.DBConn
	// var process Process

	process := struct {
		gorm.Model
		Title    string `json: "title"`
		Summary  string `json: "summary"`
		TypeID   uint   `json: "typeID"`
		CenterID uint   `json: "centerID"`
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

	dbProcess := Process{Title: process.Title, Summary: process.Summary, ProcessTypeID: process.TypeID, ProcessType: processType, Center: center, ProcessStatus: status}

	fmt.Printf("[PostProcess::dbProcess]: %v\n", dbProcess)

	db.Create(&dbProcess)
	// db.Model(&ProcessType{}).Association("Process").Append(&dbProcess)

	return c.JSON(dbProcess)
}
