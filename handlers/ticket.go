package handlers

import (
	"idolhandshake-api/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BuyTicket(c *fiber.Ctx) error {
	type BuyTicket struct {
		EventID  string `json:"event_id"`
		TicketID string `json:"ticket_id"`
		Quantity int    `json:"quantity"`
	}

	body := new(BuyTicket)

	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	if body.EventID == "" || body.TicketID == "" || body.Quantity == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Please fill all fields",
		})
	}

	userClaims := c.Locals("user").(*jwt.Token)
	claims := userClaims.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid role",
		})
	}

	objID, _ := primitive.ObjectIDFromHex(id)

	var user bson.M
	err := config.Collections.Users.FindOne(c.Context(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "User does not exist",
		})
	}

	eventID, _ := primitive.ObjectIDFromHex(body.EventID)

	var event bson.M
	err = config.Collections.Events.FindOne(c.Context(), bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Event does not exist",
		})
	}

	title := event["event_title"]

	var ticket bson.M
	for _, t := range event["tickets"].(primitive.A) {
		if t.(bson.M)["ticket_id"] == body.TicketID {
			ticket = t.(bson.M)
		}
	}

	if ticket == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Ticket does not exist",
		})
	}

	if ticketQuantity, ok := ticket["ticket_quantity"].(int32); !ok || int(ticketQuantity) < body.Quantity {
		return c.Status(400).JSON(fiber.Map{
			"message": "Not enough tickets",
		})
	}

	if qty, ok := ticket["ticket_quantity"].(int32); ok {
		ticket["ticket_quantity"] = int(qty) - body.Quantity
	}

	_, err = config.Collections.BuyTickets.InsertOne(c.Context(), bson.M{
		"event_id":        eventID,
		"ticket_id":       body.TicketID,
		"ticket_quantity": body.Quantity,
		"user_id":         objID,
		"price":           ticket["price"],
		"image":           event["event_image"],
		"ticket_name":     title,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	_, err = config.Collections.Events.UpdateOne(c.Context(), bson.M{"_id": eventID}, bson.M{
		"$set": bson.M{
			"tickets": event["tickets"],
		},
	})

	if err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Ticket bought",
	})
}

func GetTicketByUserID(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Token)
	claims := userClaims.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	role := claims["role"].(string)

	if role != "customer" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid role",
		})
	}

	objID, _ := primitive.ObjectIDFromHex(id)

	cursor, err := config.Collections.BuyTickets.Find(c.Context(), bson.M{"user_id": objID})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	var tickets []bson.M
	if err = cursor.All(c.Context(), &tickets); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	return c.Status(200).JSON(fiber.Map{
		"tickets": tickets,
	})
}