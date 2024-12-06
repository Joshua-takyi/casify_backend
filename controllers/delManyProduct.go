package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteManyProducts(ctx *gin.Context) {

	collection := Client.Database(dbName).Collection(colName)
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
