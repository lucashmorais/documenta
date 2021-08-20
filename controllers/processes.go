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

	processIDRaw := c.Query("processID")
	statusString := c.Query("statusString")
	statusIDRaw := c.Query("statusID")
	typeString := c.Query("typeString")
	typeIDRaw := c.Query("typeID")

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
