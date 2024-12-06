package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FilterProducts(ctx *gin.Context) {

	// Retrieve filter parameters from query string
	nameFilter := ctx.DefaultQuery("name", "")
	priceFilter := ctx.DefaultQuery("price", "")
	categoryFilter := ctx.DefaultQuery("category", "")
	ratingFilter := ctx.DefaultQuery("rating", "")
	discountFilter := ctx.DefaultQuery("discount", "")
	sortFilter := ctx.DefaultQuery("sort", "asc")

	// Initialize filter for database query
	filter := bson.M{}
	if nameFilter != "" {
		filter["title"] = bson.M{"$regex": nameFilter, "$options": "i"} // Case-insensitive regex search
	}
	if priceFilter != "" {
		price, err := strconv.ParseFloat(strings.TrimSpace(priceFilter), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid price format",
				"error":   err.Error(),
			})
			return
		}
		// Apply price filter based on sort order (ascending or descending)
		if sortFilter == "asc" {
			filter["price"] = bson.M{"$lte": price}
		} else if sortFilter == "desc" {
			filter["price"] = bson.M{"$gte": price}
		}
	}
	if categoryFilter != "" {
		// Split the category filter into individual categories
		categories := strings.Split(categoryFilter, ",")

		// Create a slice to hold our search conditions
		var orConditions []bson.M

		// For each category, create a search condition
		for _, category := range categories {
			// Create a condition that searches for the category (case-insensitive)
			orConditions = append(orConditions, bson.M{
				"category": bson.M{
					"$regex":   category, // Partial match
					"$options": "i",      // Case-insensitive
				},
			})
		}

		// Add the OR conditions to the filter
		filter["$or"] = orConditions
	}
	if ratingFilter != "" {
		rating, err := strconv.ParseFloat(ratingFilter, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid rating format",
				"error":   err.Error(),
			})
			return
		}
		// Use a range to account for potential type differences
		filter["rating"] = bson.M{
			"$gte": rating,
		}
	}
	if discountFilter != "" {
		discount, err := strconv.ParseFloat(discountFilter, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid discount format",
				"error":   err.Error(),
			})
			return
		}
		filter["discount"] = bson.M{
			"$gte": discount,
		}
	}

	// Define the sort order (default to ascending if not specified)
	sortOrder := bson.D{}
	if sortFilter == "asc" {
		sortOrder = bson.D{primitive.E{Key: "price", Value: 1}} // 1 for ascending
	} else if sortFilter == "desc" {
		sortOrder = bson.D{primitive.E{Key: "price", Value: -1}} // -1 for descending
	}

	// Query the MongoDB database
	collection := Client.Database(dbName).Collection(colName) // Assuming Client, dbName, and colName are defined elsewhere

	cursor, err := collection.Find(context.Background(), filter, options.Find().SetSort(sortOrder))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get products",
			"error":   err.Error(),
		})
		return
	}
	defer cursor.Close(context.Background())

	var products []model.Product
	if err := cursor.All(context.Background(), &products); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to process products",
			"error":   err.Error(),
		})
		return
	}

	// Return the filtered products in the response
	ctx.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}
