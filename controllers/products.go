package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const colName = "products"
const (
	invalidBody        = "invalid request body"
	failedToAddProduct = "Failed to add product"
	productAdded       = "Product added successfully"
	productNotFound    = "Product not found"
	productUpdated     = "Product updated successfully"
	productDeleted     = "Product deleted successfully"
	productsNotDeleted = "Failed to delete products"
)
const (
	accessControlAllowOrigin  = "Access-Control-Allow-Origin"
	accessControlAllowMethods = "Access-Control-Allow-Methods"
)

func AddProduct(ctx *gin.Context) {
	// Set CORS headers
	ctx.Header(accessControlAllowOrigin, "*")
	ctx.Header(accessControlAllowMethods, "POST")

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

func GetProducts(ctx *gin.Context) {
	ctx.Header(accessControlAllowOrigin, "*")
	ctx.Header(accessControlAllowMethods, "GET")
	collection := client.Database(dbName).Collection(colName)

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get products",
			"error":   err.Error(),
		})
		return
	}

	var products []model.Product

	if err := cursor.All(context.Background(), &products); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get products",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}

func GetById(ctx *gin.Context) {
	ctx.Header(accessControlAllowOrigin, "*")
	ctx.Header(accessControlAllowMethods, "GET")
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidBody,
			"error":   err.Error(),
		})
		return
	}
	filter := bson.M{"_id": id}

	var product model.Product

	collection := client.Database(dbName).Collection(colName)

	if err := collection.FindOne(context.Background(), filter).Decode(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": productNotFound,
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func DeleteProduct(ctx *gin.Context) {
	ctx.Header(accessControlAllowOrigin, "*")
	ctx.Header(accessControlAllowMethods, "DELETE")

	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidBody,
			"error":   err.Error(),
		})
		return
	}

	var product model.Product
	filter := bson.M{"_id": id}
	collection := client.Database(dbName).Collection(colName)

	if err := collection.FindOneAndDelete(context.Background(), filter).Decode(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": productNotFound,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": productDeleted,
	})
}

func UpdateProduct(ctx *gin.Context) {
	ctx.Header(accessControlAllowOrigin, "*")
	ctx.Header(accessControlAllowMethods, "PUT")
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidBody,
			"error":   err.Error(),
		})
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": ctx.Request.Body}

	var product model.Product

	collection := client.Database(dbName).Collection(colName)

	if err := collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": productNotFound,
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": productUpdated,
		"data":    product,
	})
}

func DeleteManyProducts(ctx *gin.Context) {
	ctx.Header(accessControlAllowOrigin, "*")
	ctx.Header(accessControlAllowMethods, "DELETE")
	collection := client.Database(dbName).Collection(colName)
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": productsNotDeleted,
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": productsNotDeleted,
	})
}
