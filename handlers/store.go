package handlers

import (
	"fmt"
	"idolhandshake-api/config"
	"idolhandshake-api/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllEvent
func GetAllEvent(c *fiber.Ctx) error {
	data, err := config.Collections.Stores.Find(c.Context(), nil)
	
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "All Ticket", "data": data})
}

// CreateEvent
func CreateEvent(c *fiber.Ctx) error {
	event := new(models.Event)
	fmt.Println(event)
	if err := c.BodyParser(event); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create ticket", "data": err})
	}

	_, err := config.Collections.Stores.InsertOne(c.Context(), event)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created Ticket", "data": event})
}
