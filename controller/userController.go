package controller

import "github.com/gofiber/fiber/v2"

func Signup() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(map[string]string{
			"message": "Signup",
		})
	}
}

// Login to add user session
func Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(map[string]string{
			"message": "Login",
		})
	}
}

// Logout to clear the session
func Logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(map[string]string{
			"message": "Logout",
		})
	}
}
