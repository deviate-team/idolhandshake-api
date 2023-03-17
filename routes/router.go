package routes

import (
	"github.com/gofiber/fiber/v2"
	"idolhandshake-api/handlers"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	auth := app.Group("/auth")

	auth.Post("/login", handlers.Login)
	auth.Post("/register", handlers.Register)
}