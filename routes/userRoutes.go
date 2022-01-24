package routes

import (
	"github.com/gofiber/fiber/v2"
	
	"github.com/braswelljr/goax/controller"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/signup", controller.Signup())
	api.Post("/login", controller.Login())
	api.Post("/logout", controller.Logout())
}
