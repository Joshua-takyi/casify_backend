package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetProducts(ctx *gin.Context) {

	collection := Client.Database(dbName).Collection(colName)

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
