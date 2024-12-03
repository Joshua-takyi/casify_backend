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
const invalidBody = "invalid request body"

func AddProduct(ctx *gin.Context) {
	inputVals := model.Product{}
	if err := ctx.ShouldBindJSON(&inputVals); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidCredentials,
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
	inputVals.Ratings = 0
	inputVals.Comments = nil
	inputVals.Discount = 0

	collection := client.Database(dbName).Collection(colName)
	_, err := collection.InsertOne(context.Background(), inputVals)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add product",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Product added successfully",
	})
}

func validateProductInput(p model.Product) error {
	// Check required fields
	if p.Title == "" || p.Description == "" || p.Price == 0 ||
		p.Images == nil || p.Details == "" || p.Color == "" {
		return errors.New("all fields are required")
	}

	// Validate price and discount
	if p.Price < 0 || p.Discount < 0 {
		return errors.New("price and discount must be non-negative")
	}

	// Validate ratings
	if p.Ratings < 0 || p.Ratings > 5 {
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

	// validate image types
	for _, img := range p.Images {
		if img != "image/jpeg" && img != "image/png" {
			return errors.New("images must be of type jpeg or png")
		}
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
