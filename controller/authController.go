package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	
	"github.com/braswelljr/goax/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
	
	"github.com/braswelljr/goax/model"
)

var (
	collection = database.OpenCollection(database.Client, "users")
	validate   = validator.New()
)

// Signup - creates and saves a new user
func Signup() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// create a new user
		user := &model.User{}
		
		// decode the request body into the user struct
		if err := ctx.BodyParser(user); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		
		// validate the user
		if err := validate.Struct(user); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		
		// hash the user's password
		password, err := HashPassword(user.Password)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		
		// check if the user already exists
		if err = collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&model.User{}); err == nil {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":      "User already exists",
				"statusCode": fiber.StatusConflict,
			})
		}
		
		// set the hashed password
		user.Id = primitive.NewObjectID()
		user.Password = password
		user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
		user.LastLogin = primitive.NewDateTimeFromTime(time.Now())
		
		// set the user's role
		user.Role = "USER"
		
		// insert the user into the database
		if _, err = collection.InsertOne(context.TODO(), user); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":      err.Error(),
				"statusCode": fiber.StatusInternalServerError,
			})
		}
		
		// return the user
		return ctx.Status(200).JSON(map[string]interface{}{
			"message":    "Signup successful",
			"payload":    user,
			"statusCode": fiber.StatusOK,
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

// HashPassword to hash the user's password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

// CheckPasswordHash to check the user's password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
