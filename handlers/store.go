package handlers

import (
	"fmt"
	"idolhandshake-api/config"
	"idolhandshake-api/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllProducts query all products
func GetAllTicket(c *fiber.Ctx) error {
	data, err := config.Collections.Stores.Find(c.Context(), nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "All Ticket", "data": data})
}

// // GetProduct query product
// func GetProduct(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	db := database.DB
// 	var product model.Product
// 	db.Find(&product, id)
// 	if product.Title == "" {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No product found with ID", "data": nil})

// 	}
// 	return c.JSON(fiber.Map{"status": "success", "message": "Product found", "data": product})
// }

// CreateTicket
func CreateTicket(c *fiber.Ctx) error {
	ticket := new(models.Ticket)
	fmt.Println(ticket)
	if err := c.BodyParser(ticket); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create ticket", "data": err})
	}

	_, err := config.Collections.Stores.InsertOne(c.Context(), ticket)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created Ticket", "data": ticket})
}

// // DeleteProduct delete product
// func DeleteProduct(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	db := database.DB

// 	var product model.Product
// 	db.First(&product, id)
// 	if product.Title == "" {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No product found with ID", "data": nil})

// 	}
// 	db.Delete(&product)
// 	return c.JSON(fiber.Map{"status": "success", "message": "Product successfully deleted", "data": nil})
// }
