package handlers

import (
	"idolhandshake-api/config"
	"idolhandshake-api/models"
	"idolhandshake-api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(c *fiber.Ctx) error {
	user := new(models.User)
	c.BodyParser(user)

	user.Password, _ = utils.HashPassword(user.Password)

	response, err := config.Collections.Users.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var insertedUser bson.M

	err = config.Collections.Users.FindOne(c.Context(), bson.M{"_id": response.InsertedID}).Decode(&insertedUser)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	token, _ := utils.GenerateToken(insertedUser["_id"].(primitive.ObjectID))

	return c.Status(201).JSON(fiber.Map{
		"access_token": token,
	})
}

func Login(c *fiber.Ctx) error {
	return c.SendString("Login")
}
