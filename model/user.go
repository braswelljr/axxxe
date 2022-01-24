package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"id"`
	Username     string             `json:"username" bson:"username"`
	Firstname    string             `json:"firstname" bson:"firstname"`
	Lastname     string             `json:"lastname" bson:"lastname"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	Password     string             `json:"password" bson:"password" validate:"required" min:"8"`
	Phone        string             `json:"phone" bson:"phone" validate:"required"`
	DateOfBirth  time.Time          `json:"date_of_birth" bson:"date_of_birth" validate:"required,date"`
	LastLogin    time.Time          `json:"last_login" bson:"last_login"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	Role         string             `json:"role" bson:"role" validate:"required, eq=ADMIN|eq=USER"`
	UserId			string 		`json:"user_id" bson:"user_id"`
}
