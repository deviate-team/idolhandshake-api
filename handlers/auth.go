package handlers

import (
	"idolhandshake-api/config"
	"idolhandshake-api/models"
	"idolhandshake-api/utils"

	"net/mail"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(c *fiber.Ctx) error {
	type Register struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body := new(Register)

	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	if body.Username == "" || body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Please fill all fields",
		})
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid email",
		})
	}

	var userExists bson.M

	err := config.Collections.Users.FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&userExists)

	if err == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "User already exists",
		})
	}

	hashedPassword, _ := utils.HashPassword(body.Password)

	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: hashedPassword,
	}

	result, err := config.Collections.Users.InsertOne(c.Context(), user)

	if err != nil {

		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	token, _ := utils.GenerateToken(result.InsertedID.(primitive.ObjectID))

	return c.Status(200).JSON(fiber.Map{
		"access_token": token,
	})
}

func Login(c *fiber.Ctx) error {
	type Login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	login := new(Login)

	c.BodyParser(login)

	if login.Email == "" || login.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Please fill all fields",
		})
	}

	if _, err := mail.ParseAddress(login.Email); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid email",
		})
	}

	var user bson.M

	err := config.Collections.Users.FindOne(c.Context(), bson.M{"email": login.Email}).Decode(&user)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	if !utils.CheckHashPassword(login.Password, user["password"].(string)) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	token, _ := utils.GenerateToken(user["_id"].(primitive.ObjectID))

	return c.Status(200).JSON(fiber.Map{
		"access_token": token,
	})
}