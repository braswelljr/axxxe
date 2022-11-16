package product

import (
	"context"
	"time"

	"github.com/braswelljr/goax/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/braswelljr/goax/database"
)

var (
	collection = database.OpenCollection(database.Client, "products")
	validate   = validator.New()
)

// GetAllProducts - Get all products
func GetAllProducts() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		return ctx.Status(200).JSON(fiber.Map{
			"status": fiber.StatusOK,
			"data":   nil,
		})
	}
}

func GetProductById(id string) (*model.Product, error) {
	// convert the id to an object id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// context
	contxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// get the user from the database
	product := &model.Product{}
	if err := collection.FindOne(contxt, bson.M{"id": oid}).Decode(product); err != nil {
		return nil, err
	}

	// return the product
	return product, nil
}

// GetProduct - get a product by id
func GetProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get the id from the of the product
		id := ctx.Params("product_id")

		// get the product
		product, err := GetProductById(id)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"status": fiber.StatusInternalServerError,
				"data":   nil,
			})
		}

		// return the product
		return ctx.Status(200).JSON(fiber.Map{
			"message":    "User found",
			"payload":    product,
			"statusCode": fiber.StatusOK,
		})
	}
}
