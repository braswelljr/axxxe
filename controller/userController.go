package controller

import (
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "log"
  "strconv"
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
    log.Println(user)
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
  if err := collection.FindOne(contxt, bson.M{"id": oid}).Decode(user); err != nil {
    return nil, err
  }
  log.Println(user)

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

    //// get all the users from the database
    var users []bson.M
    if err := result.All(contxt, &users); err != nil {
      return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error":  err,
        "status": fiber.StatusInternalServerError,
      })
    }

    // return the users
    return ctx.Status(200).JSON(map[string]interface{}{
      "message":    "Users found",
      "payload":    users[0],
      "statusCode": fiber.StatusOK,
    })
  }
}
