package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/braswelljr/axxxe/controllers/v1/authentication"
	"github.com/braswelljr/axxxe/controllers/v1/product"
	"github.com/braswelljr/axxxe/controllers/v1/user"
	"github.com/braswelljr/axxxe/middleware"
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
			auth.Post("/signup", authentication.Signup()) // Signup new users
			auth.Post("/login", authentication.Login())   // Login users
			auth.Post("/logout", authentication.Logout()) // Logout users
		}
		// Protected routes
		usr := v1.Use(middleware.Authenticate()).Group("/users")
		{
			usr.Get("/", user.GetAllUsers())                                        // Get all users
			usr.Get("/:user_id", user.GetUser())                                    // Get user by id
			usr.Patch("/:user_id", user.UpdateUser())                               // Update user by id
			usr.Patch("/:user_id/update-password", authentication.UpdatePassword()) // Update password
			usr.Patch("/:user_id/forgot-password", authentication.ForgotPassword()) // Update password
		}
	}
	// Product routes
	{
		products := v1.Group("/products")
		{
			products.Get("/", product.GetAllProducts())        // Get all products
			products.Get("/:product_id", product.GetProduct()) // Get product by id
		}
	}
}
