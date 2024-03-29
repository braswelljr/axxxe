package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User - for user params
type User struct {
	Id           primitive.ObjectID `json:"id" bson:"id"`
	Username     string             `json:"username" bson:"username" validate:"required" minlength:"3"`
	Firstname    string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname     string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	Password     string             `json:"password" bson:"password" validate:"required" min:"8"`
	Phone        string             `json:"phone,omitempty" bson:"phone,omitempty" validate:"required"`
	Gender       string             `json:"gender,omitempty" bson:"gender" validate:"required,eq=FEMALE|eq=MALE"`
	LastLogin    primitive.DateTime `json:"last_login" bson:"last_login"`
	CreatedAt    primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt    primitive.DateTime `json:"updated_at" bson:"updated_at"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	Role         string             `json:"role" bson:"role" validate:"required,eq=ADMIN|eq=USER"`
	UserId       string             `json:"user_id" bson:"user_id"`
}

// LoginDetails - email and password for user login
type LoginDetails struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

// TokenizedUserParams - used for setting the user token
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

// PasswordUpdateParams - password update params
type PasswordUpdateParams struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" bson:"password" validate:"required" min:"8"`
}
