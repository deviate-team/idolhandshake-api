package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseCollection struct {
	Users  *mongo.Collection
	Stores *mongo.Collection
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
	StoresCollection := db.Collection("stores")
	if err != nil {
		fmt.Println("Error connecting to MongoDB")
		return err
	}

	Collections = DatabaseCollection{
		Users:  usersCollection,
		Stores: StoresCollection,
	}
	fmt.Println("Connected to MongoDB")
	Client = client
	return nil
}
