package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/braswelljr/axxxe/routes"
)

func main() {
	PORT := os.Getenv("PORT")
	// check and set empty PORT
	if PORT == "" {
		PORT = "3030"
	}
	// Initialize app
	app := fiber.New()

	// Logger
	app.Use(logger.New())

	// add CORS
	app.Use(cors.New())

	// add static files
	app.Static("/", "./static")

	// index route
	app.All("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Welcome to the axxxe api",
		})
	})

	// handle routes
	routes.Routes(app)

	// Listen on port
	if err := app.Listen(":" + PORT); err != nil {
		log.Fatal("Something went wrong  ", err)
	}
}
