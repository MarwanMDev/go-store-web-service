package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Title       string             `json:"title,omitempty" validate:"required"`
	Price       string             `json:"price,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	Images      []string           `bson:"images,omitempty"`
}
