package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"github.com/braswelljr/goax/model"
)

// GetUser - gets a user by id
func GetUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get the user id from the request params
		id := ctx.Params("id")
		
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
	
	// get the user from the database
	user := &model.User{}
	if err := collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(user); err != nil {
		return nil, err
	}
	
	// return the user
	return user, nil
}
