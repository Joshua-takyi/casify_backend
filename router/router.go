package router

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/controllers"
	"github.com/joshua/casify/middleware"
)

func Router() *gin.Engine {
	r := gin.Default()

	allowOrigins := []string{"http://localhost:3000"}
	if prodOrigins := os.Getenv("ALLOWED_ORIGINS"); prodOrigins != "" {
		allowOrigins = strings.Split(prodOrigins, ",")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"origin", "content-type", "accept", "authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")

	v1.POST("/register", controllers.RegisterClient)
	v1.POST("/login", controllers.LoginClient)
	v1.POST("/addProduct", controllers.AddProduct)
	v1.POST("/addManyProducts", controllers.AddManyProducts)
	v1.GET("/getProducts", controllers.GetProducts)
	v1.GET("/validate", middleware.ValidateAuth, controllers.Validate)
	v1.GET("/filterProducts", controllers.FilterProducts)
	v1.GET("/getProduct/:id", controllers.GetById)
	v1.PUT("/updateProduct/:id", controllers.UpdateProduct)
	v1.DELETE("/deleteProduct/:id", controllers.DeleteProduct)
	v1.DELETE("/deleteProducts", controllers.DeleteManyProducts)

	return r
}
