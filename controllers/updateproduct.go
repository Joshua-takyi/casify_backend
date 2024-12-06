package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateProduct(ctx *gin.Context) {

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

	collection := Client.Database(dbName).Collection(colName)

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
