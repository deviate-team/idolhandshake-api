package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseCollection struct {
	Users *mongo.Collection
}

var Collections DatabaseCollection

var Client *mongo.Client

func ConnectDB() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(Config("MONGO_URI")))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	db := client.Database(Config("MONGO_DB"))
	usersCollection := db.Collection("users")
	if err != nil {
		return err
	}

	Collections = DatabaseCollection{
		Users: usersCollection,
	}
	Client = client
	return nil
}
