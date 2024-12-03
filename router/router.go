package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/controllers"
)

func Router() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.POST("/register", controllers.RegisterClient)
	v1.POST("/login", controllers.LoginClient)
	v1.POST("/addProduct", controllers.AddProduct)
	return r
}
