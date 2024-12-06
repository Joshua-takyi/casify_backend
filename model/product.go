package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description" binding:"required"`
	Price       float64            `json:"price" bson:"price" binding:"required"`
	Images      []string           `json:"images" bson:"images" binding:"required"`
	Discount    float64            `json:"discount,omitempty" bson:"discount,omitempty"`
	Details     ProductDetails     `json:"details" bson:"details" binding:"required"`
	Rating      float64            `json:"rating,omitempty" bson:"rating,omitempty"`
	Color       string             `json:"color,omitempty" bson:"color,omitempty"`
	Category    []string           `json:"category,omitempty" bson:"category,omitempty"`
	Comments    []string           `json:"comments,omitempty" bson:"comments,omitempty"`
	TimeStamp   TimeStamp          `json:"time_stamp,omitempty" bson:"time_stamp,omitempty"`
}

type ProductDetails struct {
	Details  []string `json:"details,omitempty" bson:"details,omitempty"`
	Features []string `json:"features,omitempty" bson:"features,omitempty"`
}
