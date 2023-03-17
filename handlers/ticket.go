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

	ticketID, _ := primitive.ObjectIDFromHex(body.TicketID)

	var ticket bson.M
	var ticketIndex int

	for i, t := range event["ticket"].([]interface{}) {
		if t.(primitive.ObjectID) == ticketID {
			ticket = t.(bson.M)
			ticketIndex = i
			break
		}
	}

	if ticket == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Ticket does not exist",
		})
	}

	if ticket["quantity"].(int) < body.Quantity {
		return c.Status(400).JSON(fiber.Map{
			"message": "Ticket is not available",
		})
	}

	ticket["quantity"] = ticket["quantity"].(int) - body.Quantity

	event["ticket"].([]interface{})[ticketIndex] = ticket

	_, err = config.Collections.Events.UpdateOne(c.Context(), bson.M{"_id": eventID}, bson.M{"$set": event})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	ticket["event_id"] = eventID
	ticket["user_id"] = objID
	ticket["quantity"] = body.Quantity
	ticket["event_image"] = event["image"]
	ticket["event_name"] = event["name"]
	ticket["price"] = ticket["price"].(int) * body.Quantity

	_, err = config.Collections.BuyTickets.InsertOne(c.Context(), ticket)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Ticket bought",
	})
}
