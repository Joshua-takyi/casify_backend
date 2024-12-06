package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/helpers"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddProduct(ctx *gin.Context) {

	inputVals := model.Product{}
	if err := ctx.ShouldBindJSON(&inputVals); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidBody,
			"error":   err.Error(),
		})
		return
	}

	if err := helpers.ValidateProductInput(inputVals); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidBody,
			"error":   err.Error(),
		})
		return
	}

	// insert default values
	inputVals.Id = primitive.NewObjectID()
	inputVals.TimeStamp.CreatedAt = time.Now()
	inputVals.TimeStamp.UpdatedAt = time.Now()

	// Set default values if not provided
	if inputVals.Comments == nil {
		inputVals.Comments = []string{}
	}
	if inputVals.Rating == 0 {
		inputVals.Rating = 0.0
	}
	if inputVals.Discount == 0 {
		inputVals.Discount = 0.0
	}

	collection := Client.Database(dbName).Collection(colName)
	_, err := collection.InsertOne(context.Background(), inputVals)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": failedToAddProduct,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":   productAdded,
		"productId": inputVals.Id,
	})
}
