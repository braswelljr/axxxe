package routes

import (
	"github.com/braswelljr/goax/controller"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// Api Routes with /api
	api := app.Group("/api")
	// User prefixed routes
	user := api.Group("/user")
	{
		user.Post("/signup", controller.Signup())
		user.Post("/login", controller.Login())
		user.Post("/logout", controller.Logout())
	}
}
