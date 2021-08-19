package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

func GetProcesses(c *fiber.Ctx) error {
	db := database.DBConn
	var processes []Process
	db.Preload("ProcessType").Find(&processes)

	return c.JSON(processes)
}

func PostProcess(c *fiber.Ctx) error {
	db := database.DBConn
	// var process Process

	process := struct {
		gorm.Model
		Title   string `json: "title"`
		Summary string `json: "summary"`
		TypeID  uint   `json: "typeID"`
	}{}

	err := c.BodyParser(&process)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"cause":   "form_decode_error",
		})
	}

	fmt.Printf("[PostProcess]: Decoded new process: %v\n", process)
	fmt.Printf("[PostProcess::process::TypeID]: %v\n", process.TypeID)

	var processType ProcessType
	db.Where(process.TypeID).Find(&processType)

	dbProcess := Process{Title: process.Title, Summary: process.Summary, ProcessTypeID: process.TypeID, ProcessType: processType}

	fmt.Printf("[PostProcess::dbProcess]: %v\n", dbProcess)

	db.Create(&dbProcess)
	// db.Model(&ProcessType{}).Association("Process").Append(&dbProcess)

	return c.JSON(dbProcess)
}
