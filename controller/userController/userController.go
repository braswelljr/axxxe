package userController

import (
  "context"
  "errors"
  "strconv"
  "time"

  "github.com/go-playground/validator/v10"
  "github.com/gofiber/fiber/v2"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"

  "github.com/braswelljr/goax/database"
  "github.com/braswelljr/goax/helper"
  "github.com/braswelljr/goax/model"
)

var (
  collection = database.OpenCollection(database.Client, "users")
  validate   = validator.New()
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
    return ctx.Status(200).JSON(fiber.Map{
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
  if err := collection.FindOne(contxt, bson.M{"id": oid}).Decode(user); err != nil {
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

    recordsPerPage, err := strconv.Atoi(ctx.Query("recordsPerPage"))
    if err != nil || recordsPerPage < 1 {
      recordsPerPage = 10
    }
    page, err := strconv.Atoi(ctx.Query("page"))
    if err != nil || page < 1 {
      page = 1
    }
    //skip := (page - 1) * recordsPerPage
    startIndex, err := strconv.Atoi(ctx.Query("startIndex"))
    if err != nil || startIndex < 1 {
      startIndex = 1
    }

    // make a match query to get all the users
    match := bson.D{{"$match", bson.D{{}}}}

    // make a group query to get the total number of users
    group := bson.D{{"$group", bson.D{
      {"_id", nil},
      {"total_count", bson.D{{"$sum", 1}}},
      {"data", bson.D{{"$push", "$$ROOT"}}},
    }}}

    // make a project query to get the users
    project := bson.D{{"$project", bson.D{
      {"_id", 0},
      {"total_count", "$total_count"},
      {"data", bson.D{{"$slice", []interface{}{"$data", startIndex, recordsPerPage}}}}},
    }}

    result, err := collection.Aggregate(contxt, mongo.Pipeline{
      match, group, project,
    })
    if err != nil {
      return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error":  err,
        "status": fiber.StatusInternalServerError,
      })
    }

    // get all the users from the database
    var users []bson.M
    if err := result.All(contxt, &users); err != nil {
      return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error":  err,
        "status": fiber.StatusInternalServerError,
      })
    }

    // return the users
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
      "message":    "Users found",
      "payload":    users,
      "statusCode": fiber.StatusOK,
    })
  }
}

// UpdateUser - updates a user
// Fields that can be updated:
//  - firstname
//  - lastname
//  - username
//  - phone
//  - gender
func UpdateUser() fiber.Handler {
  // get the user id from the request params
  return func(ctx *fiber.Ctx) error {
    // get the user id from the request params
    id := ctx.Params("user_id")

    // get the user from the database
    user, err := GetUserById(id)
    if err != nil {
      return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusInternalServerError,
      })
    }

    // set old email
    oldEmail := user.Email

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

    // check for more than one user with the same email
    if user.Email != "" {
      if err := checkEmail(user.Email); err != nil {
        return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
          "error":  err.Error(),
          "status": fiber.StatusConflict,
        })
      }
    }

    // check for more than one user with the same email again
    if err := checkEmail(user.Email); err != nil {
      user.Email = oldEmail
      return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusBadRequest,
      })
    }

    // context
    contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()

    // update the user
    user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

    // update the user in the database
    if _, err := collection.UpdateOne(contxt, bson.M{"user_id": user.UserId}, bson.M{"$set": user}); err != nil {
      return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error":  err.Error(),
        "status": fiber.StatusInternalServerError,
      })
    }

    // return the user
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
      "message":    "User updated",
      "statusCode": fiber.StatusNoContent,
    })
  }
}

// check email
func checkEmail(email string) error {
  // context
  contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
  defer cancel()

  // make a match query to get all the users
  match := bson.D{{"$match", bson.D{{"email", email}}}}

  // make a group query to get the total number of users
  group := bson.D{{"$group", bson.D{
    {"_id", nil},
    {"total_count", bson.D{{"$sum", 1}}},
  }}}

  // make a project query to get the users
  project := bson.D{{"$project", bson.D{
    {"_id", 0},
    {"total_count", "$total_count"},
  }}}

  result, err := collection.Aggregate(contxt, mongo.Pipeline{
    match, group, project,
  })
  if err != nil {
    return err
  }

  // get all the users from the database
  var users []bson.M
  if err := result.All(contxt, &users); err != nil {
    return err
  }

  // check if there is more than one user with the same email
  if len(users) > 1 {
    return errors.New("email already exists")
  }

  return nil
}
