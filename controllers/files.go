package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type FileBlob struct {
	Blob []byte
}

type UserFile struct {
	gorm.Model
	Name         string `json: name`
	FileContents string `json: fileContents`
	// FileBlob     FileBlob
	// FileBlobID   int
	Description string
	UserID      int
	User        User
	ProcessID   int
	Process     Process
}

// It seems that Go will only automatically export symbols starting with a capital letter
func NewFile(c *fiber.Ctx) error {
	db := database.DBConn
	var file UserFile
	// fileBlob := new(FileBlob)

	if !c.Is("json") {
		return c.Status(400).SendString("Request body does not include valid JSON")
	}

	if err := c.BodyParser(&file); err != nil {
		return c.Status(400).SendString("Request did not include information about the new file")
	}

	db.Create(&file)
	return c.Status(200).SendString("A new file was successfully added")
}

func NewFormFiles(c *fiber.Ctx) error {
	// Parse the multipart form:
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	// => *multipart.Form

	// Get all files from "documents" key:
	files := form.File["documents"]
	// => []*multipart.FileHeader

	// Loop through files:
	for _, file := range files {
		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
		// => "tutorial.pdf" 360641 "application/pdf"

		// Save the files to disk:
		err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename))

		// Check for errors
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFilesWithoutBlob(c *fiber.Ctx) error {
	db := database.DBConn
	var files []UserFile
	db.Preload("User").Preload("Process").Find(&files)
	// db.Preload("User").Preload("Process").Find(&files)

	return c.JSON(files)
}
