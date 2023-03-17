package handlers

import (
	"idolhandshake-api/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProfile(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Token)
	claims := userClaims.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	objID, _ := primitive.ObjectIDFromHex(id)

	var user bson.M
	err := config.Collections.Users.FindOne(c.Context(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "User does not exist",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"_id":      user["_id"],
		"username": user["username"],
		"email":    user["email"],
	})
}
