package authentication

import (
	"context"
	"time"

	"github.com/braswelljr/axxxe/controllers/v1/user"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/braswelljr/axxxe/model"
)

// HashPassword to hash the user's password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

// ComparePasswords to check the user's password
func ComparePasswords(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// UpdatePassword update users password
func UpdatePassword() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// context
		contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// password
		var password *model.PasswordUpdateParams

		// get id
		id := ctx.Params("user_id")

		// get user by id
		user, err := user.GetUserById(id)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusInternalServerError,
			})
		}

		// decode the request body into the user struct
		if err := ctx.BodyParser(&password); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusBadRequest,
			})
		}

		// validate the user
		if err := validate.Struct(password); err != nil {
			return ctx.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusExpectationFailed,
			})
		}

		// compare passwords
		if err = ComparePasswords(password.OldPassword, user.Password); err != nil {
			return ctx.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error":  "Please enter correct Password",
				"status": fiber.StatusExpectationFailed,
			})
		}

		// hash password and update user
		hash, err := HashPassword(password.NewPassword)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusBadRequest,
			})
		}
		// update the user
		user.Password = hash
		user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		// update the user in the database
		if _, err := collection.UpdateOne(contxt, bson.M{"user_id": user.UserId}, bson.M{"$set": user}); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusInternalServerError,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Password Updated Successfully",
			"status":  fiber.StatusNoContent,
		})
	}
}

// ForgotPassword - reset forgotten password
func ForgotPassword() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get user id
		id := ctx.Params("user_id")

		// get user by id
		_, err := user.GetUserById(id)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusInternalServerError,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Password Updated Successfully",
			"status":  fiber.StatusNoContent,
		})
	}
}
