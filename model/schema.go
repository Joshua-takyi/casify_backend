package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty" binding:"required,email"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty" binding:"required"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	TimeStamp TimeStamp          `json:"time_stamp,omitempty" bson:"time_stamp,omitempty"`
}

type TimeStamp struct {
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
