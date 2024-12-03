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
    Details     string             `json:"details" bson:"details" binding:"required"`
    Ratings     float64            `json:"ratings,omitempty" bson:"ratings,omitempty"`
    Color       string             `json:"color" bson:"color" binding:"required"`
    Comments    []string           `json:"comments,omitempty" bson:"comments,omitempty"`
    TimeStamp   TimeStamp          `json:"time_stamp" bson:"time_stamp"`
}

type TimeStamp struct {
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}



type LoginRequest struct {
    Email    string `json:"email,omitempty" bson:"email,omitempty" binding:"required,email"`
    Password string `json:"password,omitempty"  bson:"password,omitempty" binding:"required"`
}