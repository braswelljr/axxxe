package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/braswelljr/axxxe/helper"
)

// Authenticate is a middleware that checks if the user is authenticated
func Authenticate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get the token from the request
		token := ctx.Get("token")
		if token == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"status":  fiber.StatusUnauthorized,
			})
		}
		// check if the token is valid
		claims, err := helper.ValidateToken(token)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"status":  fiber.StatusUnauthorized,
			})
		}
		// set the claims to the context
		ctx.Locals("email", claims.User.Email)
		ctx.Locals("username", claims.User.Username)
		ctx.Locals("firstname", claims.User.Firstname)
		ctx.Locals("lastname", claims.User.Lastname)
		ctx.Locals("phone", claims.User.Phone)
		ctx.Locals("gender", claims.User.Gender)
		ctx.Locals("role", claims.User.Role)
		ctx.Locals("user_id", claims.User.UserId)
		return ctx.Next()
	}
}
