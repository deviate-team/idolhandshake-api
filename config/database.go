package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseCollection struct {
	Users      *mongo.Collection
	Events     *mongo.Collection
	BuyTickets *mongo.Collection
}

var Collections DatabaseCollection

var Client *mongo.Client

func ConnectDB() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(Config("MONGO_URI")))
	if err != nil {
		fmt.Println("Error connecting to MongoDB")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	db := client.Database(Config("MONGO_DB"))

	usersCollection := db.Collection("users")
	eventsCollection := db.Collection("events")
	buyTicketsCollection := db.Collection("buyTickets")

	if err != nil {
		fmt.Println("Error connecting to MongoDB")
		return err
	}

	Collections = DatabaseCollection{
		Users:      usersCollection,
		Events:     eventsCollection,
		BuyTickets: buyTicketsCollection,
	}
	fmt.Println("Connected to MongoDB")
	Client = client
	return nil
}
