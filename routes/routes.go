package routes

import (
  "github.com/gofiber/fiber/v2"

  "github.com/braswelljr/goax/controller"
  "github.com/braswelljr/goax/controller/authenticationController"
  "github.com/braswelljr/goax/controller/userController"
  "github.com/braswelljr/goax/middleware"
)

// Routes handles application routes.
// - APIs are prefixed with `api`.
// - Versions are prefixed with `v(number)`. Example `v1`, `v2`.
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
      auth.Post("/signup", authenticationController.Signup()) // Signup new users
      auth.Post("/login", authenticationController.Login())   // Login users
      auth.Post("/logout", authenticationController.Logout()) // Logout users
    }
    // Protected routes
    user := v1.Use(middleware.Authenticate()).Group("/users")
    {
      user.Get("/", userController.GetAllUsers())                                        // Get all users
      user.Get("/:user_id", userController.GetUser())                                    // Get user by id
      user.Patch("/:user_id", userController.UpdateUser())                               // Update user by id
      user.Patch("/:user_id/update-password", authenticationController.UpdatePassword()) // Update password
      user.Patch("/:user_id/forgot-password", authenticationController.ForgotPassword()) // Update password
    }
  }
  // Product routes
  {
    product := v1.Group("/products")
    {
      product.Get("/", controller.GetAllProducts())        // Get all products
      product.Get("/:product_id", controller.GetProduct()) // Get product by id
    }
  }
}
