package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/helpers"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddManyProducts(ctx *gin.Context) {
	// Step 1: Parse and bind the incoming JSON into the inputVals slice
	inputVals := make([]model.Product, 0)
	if err := ctx.ShouldBindJSON(&inputVals); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid body",
			"error":   err.Error(),
		})
		return
	}

	// Step 2: Validate each product input
	for _, inputVal := range inputVals {
		if err := helpers.ValidateProductInput(inputVal); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid body",
				"error":   err.Error(),
			})
			return
		}
	}

	// Step 3: Set default values where needed
	for i := range inputVals {
		inputVals[i].Id = primitive.NewObjectID()     // Generate a new ObjectId
		inputVals[i].TimeStamp.CreatedAt = time.Now() // Set CreatedAt
		inputVals[i].TimeStamp.UpdatedAt = time.Now() // Set UpdatedAt

		// Default empty slice for comments if not provided
		if inputVals[i].Comments == nil {
			inputVals[i].Comments = []string{}
		}
		// Set default rating if it's zero (if this is required)
		if inputVals[i].Rating == 0.0 {
			inputVals[i].Rating = 0.0
		}
		// Set default discount if it's zero
		if inputVals[i].Discount == 0.0 {
			inputVals[i].Discount = 0.0
		}
	}

	// Step 4: Ensure that the MongoDB client and collection are valid
	if Client == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to connect to the database",
		})
		return
	}

	collection := Client.Database(dbName).Collection(colName)
	if collection == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to access the collection",
		})
		return
	}
	doc := make([]interface{}, len(inputVals))

	// create an interface
	for i, d := range inputVals {
		doc[i] = d
	}

	// Step 5: Attempt to insert products into MongoDB
	_, err := collection.InsertMany(context.Background(), doc)
	if err != nil {
		// Log the error for further inspection
		fmt.Println("Failed to insert products:", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add products",
			"error":   err.Error(),
		})
		return
	}

	// Step 6: Successfully added products
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Products added successfully",
	})
}