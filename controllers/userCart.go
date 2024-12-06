package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	Id primitive.ObjectID `json:"id"`
}

func UserCart(ctx *gin.Context) {
	// Get user from context
	user, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	// Cast user to UserResponse
	userDetails, ok := user.(UserResponse)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve user details",
		})
		return
	}

	userId := userDetails.Id // This is now the correct primitive.ObjectID

	// Bind JSON input
	cartInput := model.UserCart{}
	if err := ctx.ShouldBindJSON(&cartInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Insert the cart into a single collection
	colName := userId.Hex()
	collection := Client.Database(dbName).Collection(colName)
	
	newCart := model.UserCart{
		User:     userId,
		Items:    cartInput.Items,
		Subtotal: cartInput.Subtotal,
		Tax:      cartInput.Tax,
		Shipping: cartInput.Shipping,
		Total:    cartInput.Total,
		Status:   cartInput.Status,
		TimeStamp: model.TimeStamp{
			CreatedAt: cartInput.TimeStamp.CreatedAt,
			UpdatedAt: cartInput.TimeStamp.UpdatedAt,
		},
		Currency:        cartInput.Currency,
		Discounts:       cartInput.Discounts,
		ShippingAddress: cartInput.ShippingAddress,
		PaymentMethod:   cartInput.PaymentMethod,
	}

	result, err := collection.InsertOne(context.Background(), newCart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add product to cart",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cart updated successfully",
		"cartId":  result.InsertedID,
	})
}
