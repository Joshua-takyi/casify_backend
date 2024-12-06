package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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



func handleProductError(ctx *gin.Context, err error) {
	switch {
	case strings.Contains(err.Error(), "invalid product ID"):
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidBody,
			"error":   err.Error(),
		})
	case strings.Contains(err.Error(), "product not found"):
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": productNotFound,
			"error":   err.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "unexpected error",
			"error":   err.Error(),
		})
	}
}




