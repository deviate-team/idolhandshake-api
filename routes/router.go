package routes

import (
	"idolhandshake-api/config"
	"idolhandshake-api/handlers"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	auth := app.Group("/auth")

	auth.Post("/login", handlers.Login)
	auth.Post("/register", handlers.Register)

	user := app.Group("/user")

	user.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Config("JWT_SECRET")),
	}))
	user.Get("/profile", handlers.GetProfile)

	store := app.Group("/store")
	store.Post("/add-ticket", handlers.CreateTicket)
	store.Post("/all-ticket", handlers.GetAllTicket)
}
