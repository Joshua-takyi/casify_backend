package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteProduct(ctx *gin.Context) {

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
	collection := Client.Database(dbName).Collection(colName)

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
