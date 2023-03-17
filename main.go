package main

import (
	"context"
	"log"

	"idolhandshake-api/config"
	"idolhandshake-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if err := config.ConnectDB(); err != nil {
		panic(err)
	}

	defer config.Client.Disconnect(context.Background())

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + config.Config("PORT")))
}
