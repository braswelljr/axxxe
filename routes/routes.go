package routes

import (
  "github.com/gofiber/fiber/v2"

  "github.com/braswelljr/goax/controller"
  "github.com/braswelljr/goax/middleware"
)

func Routes(app *fiber.App) {
  // API Routes with /api
  api := app.Group("/api")
  // Versioning
  // Version 1 (prefix - v1)
  v1 := api.Group("/v1")
  // User prefixed routes
  // Authentication
  {
    auth := v1.Group("/users")
    {
      auth.Post("/signup", controller.Signup())
      auth.Post("/login", controller.Login())
      auth.Post("/logout", controller.Logout())
    }
    // Protected routes
    user := v1.Use(middleware.Authenticate()).Group("/users")
    {
      user.Get("/", controller.GetAllUsers()) // Get all users
      user.Get("/:user_id", controller.GetUser()) // Get user by id
      user.Patch("/:user_id", controller.UpdateUser()) // Update user by id
    }
  }
}
