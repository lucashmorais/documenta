package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
)

type FileBlob struct {
	Blob []byte
}

type UserFile struct {
	gorm.Model
	Name        string `json: name`
	UUID        string `json: uuid`
	ContentType string `json: content-type`
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
	db := database.DBConn
	var dbFile UserFile

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["documents"]

	for _, file := range files {
		uuid := uuid.New().String()
		fmt.Println("[NewFormFiles]: ", file.Filename, file.Size, file.Header["Content-Type"][0])

		err := c.SaveFile(file, fmt.Sprintf("./user_data/%s", file.Filename))

		if err != nil {
			return err
		}

		dbFile.Name = file.Filename
		dbFile.UUID = uuid
		dbFile.ContentType = file.Header["Content-Type"][0]
		db.Create(&dbFile)
	}
	return c.JSON(map[string]string{"status": "success"})
}

func GetFilesWithoutBlob(c *fiber.Ctx) error {
	db := database.DBConn
	var files []UserFile
	db.Preload("User").Preload("Process").Find(&files)
	// db.Preload("User").Preload("Process").Find(&files)

	return c.JSON(files)
}
