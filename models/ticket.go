package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID             primitive.ObjectID `json:"id"  bson:"_id,omitempty"`
	Image          string             `json:"image" bson:"image"`
	TicketID       string             `json:"ticket_id" bson:"ticket_id"`
	EventID        primitive.ObjectID `json:"event_id" bson:"event_id"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	Price          float64            `json:"price" bson:"price"`
	TicketName     string             `json:"ticket_name" bson:"ticket_name"`
	TicketQuantity int                `json:"ticket_quantity" bson:"ticket_quantity"`
}
