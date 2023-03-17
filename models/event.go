package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID               primitive.ObjectID `json:"id"  bson:"_id,omitempty"`
	EventTitle       string             `json:"event_title" bson:"event_title"`
	EventDescription string             `json:"event_description" bson:"event_description"`
	EventDate        string             `json:"event_date" bson:"event_date"`
	EventTime        string             `json:"event_time" bson:"event_time"`
	EventLocation    string             `json:"event_location" bson:"event_location"`
	EventImage       string             `json:"event_image" bson:"event_image"`
	Tickets          []Ticket           `json:"tickets" bson:"tickets"`
	OrganizerID      primitive.ObjectID `json:"organizer_id" bson:"organizer_id"`
}
