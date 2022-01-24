package main

import (
	"github.com/braswelljr/goax/routes"
	"log"
	"os"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	PORT := os.Getenv("PORT")
	// check and set empty PORT
	if PORT == "" {
		PORT = "8080"
	}
	// Initialize app
	app := fiber.New()
	
	// Logger
	app.Use(logger.New())
	
	// add CORS
	app.Use(cors.New())
	
	// add static files
	app.Static("/", "./static")
	
	// handle routes
	routes.UserRoutes(app)
	
	// Listen on port
	if err := app.Listen(":" + PORT); err != nil {
		log.Fatal("Something went wrong  ", err)
	}
}
