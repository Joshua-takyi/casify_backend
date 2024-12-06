package model

type LoginRequest struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty" binding:"required,email"`
	Password string `json:"password,omitempty"  bson:"password,omitempty" binding:"required"`
}
