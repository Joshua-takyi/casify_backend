package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const colName = "products"
const (
	invalidBody        = "invalid request body"
	failedToAddProduct = "Failed to add product"
	productAdded       = "Product added successfully"
)

func AddProduct(ctx *gin.Context) {
	// Set CORS headers
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST")

	inputVals := model.Product{}
	if err := ctx.ShouldBindJSON(&inputVals); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidBody,
			"error":   err.Error(),
		})
		return
	}

	if err := validateProductInput(inputVals); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidBody,
			"error":   err.Error(),
		})
		return
	}

	// Proceed with product addition logic

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

	collection := client.Database(dbName).Collection(colName)
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

func validateProductInput(p model.Product) error {
	// Check required fields
	if p.Title == "" || p.Description == "" || p.Price == 0 ||
		p.Images == nil || p.Details.Details == nil || p.Color == "" {
		return errors.New("all fields are required")
	}

	// Validate price and discount
	if p.Price < 0 || p.Discount < 0 {
		return errors.New("price and discount must be non-negative")
	}

	// Validate ratings
	if p.Rating < 0 || p.Rating > 5 {
		return errors.New("ratings must be between 0 and 5")
	}

	// Validate comments
	if err := validateStringSlice(p.Comments, "comments"); err != nil {
		return err
	}

	// Validate images
	if err := validateStringSlice(p.Images, "images"); err != nil {
		return err
	}

	return nil
}

func validateStringSlice(slice []string, fieldName string) error {
	if len(slice) > 0 {
		for _, item := range slice {
			if item == "" {
				return fmt.Errorf("%s cannot contain empty strings", fieldName)
			}
		}
	}
	return nil
}
