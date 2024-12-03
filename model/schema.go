package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty" binding:"required"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty" binding:"required"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty" binding:"required,email"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty" binding:"required"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty" binding:"required,oneof=admin user"`
	TimeStamp TimeStamp          `json:"time_stamp,omitempty" bson:"time_stamp,omitempty"`
}

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
	Comments    []string           `json:"comments,omitempty" bson:"comments,omitempty"`
	TimeStamp   TimeStamp          `json:"time_stamp,omitempty" bson:"time_stamp,omitempty"`
}

type ProductDetails struct {
	ProductId primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Details   []string           `json:"details,omitempty" bson:"details,omitempty"`
	Features  []string           `json:"features,omitempty" bson:"features,omitempty"`
}

type TimeStamp struct {
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty" binding:"required,email"`
	Password string `json:"password,omitempty"  bson:"password,omitempty" binding:"required"`
}
