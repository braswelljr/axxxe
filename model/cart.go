package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Cart is a struct
type Cart struct {
	Id       primitive.ObjectID `json:"id" bson:"id"`
	UserID   int64              `json:"user_id" bson:"user_id"`
	Products []Product          `json:"product" bson:"products"`
	Price    float32            `json:"product_price" bson:"product_price"`
	Quantity int64              `json:"quantity" bson:"quantity"`
}
