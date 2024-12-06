package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetById(ctx *gin.Context) {

	product, err := findProductById(ctx)
	if err != nil {
		handleProductError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func findProductById(ctx *gin.Context) (*model.Product, error) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	filter := bson.M{"_id": id}
	collection := Client.Database(dbName).Collection(colName)

	var product model.Product
	if err := collection.FindOne(context.Background(), filter).Decode(&product); err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return &product, nil
}
