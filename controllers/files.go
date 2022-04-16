package controllers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/documenta/database"
	"github.com/lucashmorais/documenta/utils"
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

// Function that adds a new item to the File gorm table based on the data provided in the POST request as a JSON
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
	var fileIDs []int

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["documents"]
	processID, err := strconv.Atoi(c.FormValue("processID"))
	if err != nil {
		processID = 0
	}

	for _, file := range files {
		uuid := uuid.New().String()
		fmt.Println("[NewFormFiles]: ", file.Filename, file.Size, file.Header["Content-Type"][0])

		base_path := utils.GetServerPath()
		err := c.SaveFile(file, fmt.Sprintf(base_path + "/user_data/.%s", uuid))

		if err != nil {
			return err
		}

		dbFile.Name = file.Filename
		dbFile.UUID = uuid
		dbFile.ContentType = file.Header["Content-Type"][0]
		dbFile.UserID = RetrieveUserID(c)
		dbFile.ProcessID = processID
		db.Create(&dbFile)
		fileIDs = append(fileIDs, int(dbFile.ID))
	}

	fmt.Println(fileIDs)
	return c.JSON(fileIDs)
	// return c.JSON(map[string]string{"status": "success"})
}

func GetFilesWithoutBlob(c *fiber.Ctx) error {
	db := database.DBConn
	var files []UserFile
	processIDRaw := c.Query("processID")

	driver := db.Preload("User").Preload("Process")
	if processIDRaw != "" {
		i, err := strconv.Atoi(processIDRaw)
		if err == nil {
			driver = driver.Where("process_id = ?", i)
		}
	}

	driver.Find(&files)

	return c.JSON(files)
}

func GetFile(c *fiber.Ctx) error {
	db := database.DBConn
	var file UserFile
	id := c.Params("id")
	db.Find(&file, id)

	base_path := utils.GetServerPath()

	// return c.Download("./user_data/."+file.UUID, file.Name)
	return c.SendFile(base_path + "/user_data/." + file.UUID)
}

func DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var file UserFile
	db.First(&file, id)

	base_path := utils.GetServerPath()

	uuid := file.UUID
	name := file.Name

	if uuid == "" {
		return c.Status(500).SendString("File was not found")
	}

	os.Remove(base_path + "/user_data/." + uuid)
	db.Delete(&file)

	return c.SendString("The file " + name + " was successfully deleted")
}
