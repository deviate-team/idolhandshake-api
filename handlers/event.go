package handlers

import (
	"idolhandshake-api/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEvent(c *fiber.Ctx) error {
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

	if user["role"] != "organizer" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid role",
		})
	}

	type Ticket struct {
		Name     string `json:"name"`
		Price    int    `json:"price"`
		Quantity int    `json:"quantity"`
	}

	type Event struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Location    string   `json:"location"`
		Date        string   `json:"date"`
		Time        string   `json:"time"`
		Image       string   `json:"image"`
		Ticket      []Ticket `json:"ticket"`
	}

	body := new(Event)

	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	if body.Name == "" || body.Description == "" || body.Location == "" || body.Date == "" || body.Time == "" || body.Image == "" || body.Ticket == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Please fill all fields",
		})
	}

	var tickets []bson.M

	for _, ticket := range body.Ticket {
		tickets = append(tickets, bson.M{
			"ticket_name":     ticket.Name,
			"ticket_quantity": ticket.Quantity,
			"price":           ticket.Price,
			"ticket_id":       (uuid.New()).String(),
		})
	}

	event, err := config.Collections.Events.InsertOne(c.Context(), bson.M{
		"event_title":       body.Name,
		"event_description": body.Description,
		"event_location":    body.Location,
		"event_date":        body.Date,
		"event_time":        body.Time,
		"event_image":       body.Image,
		"tickets":           tickets,
		"organizer_id":      objID,
	})

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Event created",
		"event":   event,
	})
}

func GetEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")

	objID, _ := primitive.ObjectIDFromHex(eventID)

	var event bson.M
	err := config.Collections.Events.FindOne(c.Context(), bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Event does not exist",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"event": event,
	})
}

func GetEvents(c *fiber.Ctx) error {
	var events []bson.M
	cursor, err := config.Collections.Events.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	if err = cursor.All(c.Context(), &events); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"events": events,
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")

	objID, _ := primitive.ObjectIDFromHex(eventID)

	var event bson.M
	err := config.Collections.Events.FindOne(c.Context(), bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Event does not exist",
		})
	}

	_, err = config.Collections.Events.DeleteOne(c.Context(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	

	return c.Status(200).JSON(fiber.Map{
		"message": "Event deleted",
	})
}
