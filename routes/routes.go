package routes

import (
	"github.com/braswelljr/goax/controller"
	"github.com/braswelljr/goax/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// Api Routes with /api
	api := app.Group("/api")
	// User prefixed routes
	// Authentication
	auth := api.Group("/users")
	{
		auth.Post("/signup", controller.Signup())
		auth.Post("/login", controller.Login())
		auth.Post("/logout", controller.Logout())
	}
	// Protected routes
	user := api.Use(middleware.Authenticate()).Group("/users")
	{
		user.Post("/:user_id", controller.GetUser())
	}
}
