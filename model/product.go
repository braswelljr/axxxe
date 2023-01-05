package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product - for product params
type Product struct {
	Id           primitive.ObjectID `json:"id" bson:"id"`
	Image        string             `json:"image" bson:"image"`
	Name         string             `json:"name" bson:"name"`
	Type         string             `json:"type" bson:"type"`
	Description  string             `json:"description" bson:"description"`
	Price        float64            `json:"price" bson:"price"`
	Quantity     int                `json:"quantity" bson:"quantity"`
	Availability bool               `json:"availability" bson:"availability" validate:"required"`
	ProductId    string             `json:"product_id" bson:"product_id"`
}
