package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID          primitive.ObjectID `json:"id"  bson:"_id,omitempty"`
	Description     string             `json:"description" bson:"description"`
	Maxticket  int             `json:"maxticket" bson:"maxticket"`
}
