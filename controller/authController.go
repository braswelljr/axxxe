package controller

import (
	"context"
	"github.com/braswelljr/goax/helper"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/braswelljr/goax/database"
	"github.com/braswelljr/goax/model"
)

var (
	collection = database.OpenCollection(database.Client, "users")
	validate   = validator.New()
	// context
	contxt, cancel = context.WithTimeout(context.Background(), 100*time.Second)
)

// Signup - creates and saves a new user
func Signup() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// create a new user
		user := &model.User{}

		// decode the request body into the user struct
		if err := ctx.BodyParser(&user); err != nil {
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
				"error":  err.Error(),
				"status": fiber.StatusInternalServerError,
			})
		}

		// check if the user already exists
		err = collection.FindOne(contxt, bson.M{"email": user.Email}).Decode(&model.User{})
		defer cancel()
		if err == nil {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":  "User already exists",
				"status": fiber.StatusConflict,
			})
		}

		// set the hashed password
		user.Id = primitive.NewObjectID()
		user.Password = password
		user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
		user.LastLogin = primitive.NewDateTimeFromTime(time.Now())
		user.UserId = user.Id.Hex()

		// params to be tokenized
		tokenParams := &model.TokenizedUserParams{
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Password,
			Phone:     user.Phone,
			Gender:    user.Gender,
			Role:      user.Role,
			UserId:    user.UserId,
		}

		// get tokens
		token, refreshToken, err := helper.GetAllTokens(*tokenParams)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusInternalServerError,
			})
		}

		// set user token
		user.Token = token
		user.RefreshToken = refreshToken

		// insert the user into the database
		_, err = collection.InsertOne(contxt, user)
		defer cancel()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusInternalServerError,
			})
		}

		// return the user
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Signup successful",
			"payload": fiber.Map{
				"user_id":      user.UserId,
				"token":        token,
				"refreshToken": refreshToken,
			},
			"status": fiber.StatusOK,
		})
	}
}

// Login to add user session
func Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// create a new user
		user := &model.User{}
		foundUser := &model.User{}

		// decode the request body into the user struct
		if err := ctx.BodyParser(&user); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusBadRequest,
			})
		}

		// validate the user
		if err := validate.Struct(user); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  err.Error(),
				"status": fiber.StatusBadRequest,
			})
		}

		// check if the user exists
		err := collection.FindOne(contxt, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":  "Invalid email",
				"status": fiber.StatusNotFound,
			})
		}

		// check if the password is correct
		_,err = ComparePasswords(user.Password, foundUser.Password)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":  "Invalid Credentials",
				"status": fiber.StatusUnauthorized,
			})
		}

		// params to be tokenized
		tokenParams := &model.TokenizedUserParams{
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Password,
			Phone:     user.Phone,
			Gender:    user.Gender,
			Role:      user.Role,
			UserId:    user.UserId,
		}

		// get tokens
		token, refreshToken, err := helper.GetAllTokens(*tokenParams)

		// set user tokens
		user.Token = token
		user.RefreshToken = refreshToken

		// return the user
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Login successful",
			"payload": fiber.Map{
				"user_id":      user.UserId,
				"token":        token,
				"refreshToken": refreshToken,
			},
			"status": fiber.StatusOK,
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

// ComparePasswords to check the user's password
func ComparePasswords(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err.Error()
	}
	return true, nil
}
