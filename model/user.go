package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"id"`
	Username     string             `json:"username" bson:"username" validate:"required" minlength:"3"`
	Firstname    string             `json:"firstname" bson:"firstname"`
	Lastname     string             `json:"lastname" bson:"lastname"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	Password     string             `json:"password" bson:"password" validate:"required" min:"8"`
	Phone        string             `json:"phone" bson:"phone" validate:"required"`
	Gender       string             `json:"gender,omitempty" bson:"gender" validate:"required,eq=FEMALE|eq=MALE"`
	LastLogin    primitive.DateTime `json:"last_login" bson:"last_login"`
	CreatedAt    primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt    primitive.DateTime `json:"updated_at" bson:"updated_at"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	Role         string             `json:"role" bson:"role" validate:"required,eq=ADMIN|eq=USER"`
	UserId       string             `json:"user_id" bson:"user_id"`
}

type TokenizedUserParams struct {
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Phone     string
	Gender    string
	Role      string
	UserId    string
}
