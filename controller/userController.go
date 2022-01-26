package controller

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/braswelljr/goax/helper"
	"github.com/braswelljr/goax/model"
)

// GetUser - gets a user by id
func GetUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get the user id from the request params
		id := ctx.Params("user_id")

		// get user with admin role
		if err := helper.MatchUserTypeToUID(ctx, id); err != nil {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusForbidden,
			})
		}

		// get the user from the database
		user, err := GetUserById(id)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusInternalServerError,
			})
		}

		// return the user
		return ctx.Status(200).JSON(map[string]interface{}{
			"message":    "User found",
			"payload":    user,
			"statusCode": fiber.StatusOK,
		})
	}
}

// GetUserById - gets a user by id
func GetUserById(id string) (*model.User, error) {
	// convert the id to an object id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// context
	contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// get the user from the database
	user := &model.User{}
	if err := collection.FindOne(contxt, bson.M{"user_id": oid}).Decode(user); err != nil {
		return nil, err
	}

	// return the user
	return user, nil
}

// GetAllUsers fetches all the users from the database - admin only
func GetAllUsers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Check user with admin role
		if err := helper.CheckUserType(ctx, "ADMIN"); err != nil {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusForbidden,
			})
		}

		// context
		contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// get all the users from the database
		var users []*model.User
		cursor, _ := collection.Find(contxt, bson.M{})
		if err := cursor.All(contxt, &users); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  err,
				"status": fiber.StatusInternalServerError,
			})
		}

		// return the users
		return ctx.Status(200).JSON(map[string]interface{}{
			"message":    "Users found",
			"payload":    users,
			"statusCode": fiber.StatusOK,
		})
	}
}
