package products

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/go_fiber/database"
)

type Product struct {
	gorm.Model
	Name     string `json: "name"`
	Price    int    `json: "price"`
	Quantity int    `json: "quantity"`
}

func GetProducts(c *fiber.Ctx) error {
	db := database.DBConn
	var products []Product
	db.Find(&products)
	return c.JSON(products)
	// return c.SendString("All Products")
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var product Product
	db.Find(&product, id)
	return c.JSON(product)
	// return c.SendString("Single product")
}

func NewProduct(c *fiber.Ctx) error {
	db := database.DBConn
	var product Product

	// This does not ensure that all Product fields are received nor that
	// requests containing extraneous information is discarded
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).SendString("Request did not include information about the new product")
	}

	db.Create(&product)
	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var product Product
	db.First(&product, id)
	if product.Name == "" {
		return c.Status(500).SendString("Product was not found")
	}
	name := product.Name
	db.Delete(&product)
	return c.SendString("The product " + name + " was successfully deleted")
}

func IncreaseQuantity(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var product Product
	db.First(&product, id)
	if product.Name == "" {
		return c.Status(500).SendString("Product was not found")
	}
	product.Quantity = product.Quantity - 1
	db.Save(&product)
	return c.JSON(product)
}

func DecreaseQuantity(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var product Product
	db.First(&product, id)
	if product.Name == "" {
		return c.Status(500).SendString("Product was not found")
	}
	product.Quantity = product.Quantity - 1
	db.Save(&product)
	return c.JSON(product)
}
