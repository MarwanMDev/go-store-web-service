package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Email    string             `bson:"email,omitempty" validate:"required"`
	Title    string             `json:"title,omitempty" validate:"required"`
	Password string             `bson:"password,omitempty" validate:"required"`
	Role     string             `bson:"role,omitempty" validate:"required"`
	Avatar   string             `bson:"avatar,omitempty"`
}
