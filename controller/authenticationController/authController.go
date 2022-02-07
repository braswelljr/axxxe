package authenticationController

import (
  "context"
  "time"

  "github.com/go-playground/validator/v10"
  "github.com/gofiber/fiber/v2"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"

  "github.com/braswelljr/goax/database"
  "github.com/braswelljr/goax/helper"
  "github.com/braswelljr/goax/model"
)

var (
  collection = database.OpenCollection(database.Client, "users")
  validate   = validator.New()
)

// Signup - creates and saves a new user
func Signup() fiber.Handler {
  return func(ctx *fiber.Ctx) error {
    // context
    contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()
    // create a new user
    user := &model.User{}

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

    // hash the user's password
    password, err := HashPassword(user.Password)
    if err != nil {
      return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusInternalServerError,
      })
    }

    // set the hashed password
    user.Id = primitive.NewObjectID()
    user.Password = password
    user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
    user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
    user.LastLogin = primitive.NewDateTimeFromTime(time.Now())
    user.UserId = user.Id.Hex()

    // check if the user already exists
    err = collection.FindOne(contxt, bson.M{"email": user.Email}).Decode(&model.User{})
    if err == nil {
      return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
        "error":  "User already exists",
        "status": fiber.StatusConflict,
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
// Uses email and password to login
func Login() fiber.Handler {
  return func(ctx *fiber.Ctx) error {
    // context
    contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()
    // get user params for login
    var user *model.LoginDetails
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
      return ctx.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusExpectationFailed,
      })
    }

    // check if the user exists
    err := collection.FindOne(contxt, bson.M{"email": user.Email}).Decode(&foundUser)
    if err != nil {
      return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "error":  "Invalid email",
        "status": fiber.StatusNotFound,
      })
    }

    // check if the password is correct
    err = ComparePasswords(user.Password, foundUser.Password)
    if err != nil {
      return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error":  "Invalid Credentials",
        "status": fiber.StatusUnauthorized,
      })
    }

    // params to be tokenized
    tokenParams := &model.TokenizedUserParams{
      Username:  foundUser.Username,
      Firstname: foundUser.Firstname,
      Lastname:  foundUser.Lastname,
      Email:     foundUser.Password,
      Phone:     foundUser.Phone,
      Gender:    foundUser.Gender,
      Role:      foundUser.Role,
      UserId:    foundUser.UserId,
    }

    // get tokens
    token, refreshToken, _ := helper.GetAllTokens(*tokenParams)

    // update token in database
    _, err = collection.UpdateOne(
      contxt,
      bson.M{"user_id": foundUser.UserId},
      bson.M{
        "$set": bson.M{
          "token":         token,
          "refresh_token": refreshToken,
          "updated_at":    primitive.NewDateTimeFromTime(time.Now()),
          "last_login":    primitive.NewDateTimeFromTime(time.Now()),
        },
      },
    )
    if err != nil {
      return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusInternalServerError,
      })
    }

    // return the user
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
      "message": "Login successful",
      "payload": fiber.Map{
        "user_id":      foundUser.UserId,
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
    // context
    contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()
    // get user params for logout
    user := struct {
      UserId string `json:"user_id" bson:"user_id" validate:"required"`
    }{}

    // decode the request body into the user struct
    if err := ctx.BodyParser(&user); err != nil {
      return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusBadRequest,
      })
    }

    // validate the user
    if err := validate.Struct(user); err != nil {
      return ctx.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusExpectationFailed,
      })
    }

    // update token in database
    _, err := collection.UpdateOne(
      contxt,
      bson.M{"user_id": user.UserId},
      bson.M{
        "$set": bson.M{
          "token":         "",
          "refresh_token": "",
          "updated_at":    primitive.NewDateTimeFromTime(time.Now()),
        },
      },
    )

    if err != nil {
      return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusInternalServerError,
      })
    }

    // return the user
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
      "message": "Logout successful",
      "payload": fiber.Map{
        "user_id": user.UserId,
      },
      "status": fiber.StatusOK,
    })
  }
}
